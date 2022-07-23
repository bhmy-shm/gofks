package gofk

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type HttpSrv interface {
	SetAddr(srv *Server)
	Build(srv *Server)
}

type handle func(http.ResponseWriter, *http.Request)

type router struct {
	path    string
	handler handle
}

type Server struct {
	handle map[string][]*router      //路由map集合	string = 地址，router = 路由
	srv    map[string]*http.ServeMux //对应的http监听服务  string = 地址, ServeMux = http服务
	proxy  *httputil.ReverseProxy    //反向代理工具

	director func(req *http.Request)        //请求源
	modify   func(res *http.Response) error //代理的服务
}

func NewHttp() *Server {
	return &Server{srv: make(map[string]*http.ServeMux, 0), handle: make(map[string][]*router, 0)}
}

func (s *Server) Mount(srvs ...HttpSrv) *Server {
	for _, srv := range srvs {
		srv.SetAddr(s)
		srv.Build(s)
	}
	return s
}

func (s *Server) SetServe(address string) {
	if s.srv == nil {
		s.srv = make(map[string]*http.ServeMux, 0)
	}
	s.srv[address] = http.NewServeMux()
	if s.handle == nil {
		s.handle = make(map[string][]*router, 0)
	}
	s.handle[address] = make([]*router, 0)
}

func (s *Server) HandleFunc(address string, path string, fn handle) {
	if s.handle == nil {
		s.handle = make(map[string][]*router, 0)
	}

	//通过服务端地址找路由
	if _, found := s.handle[address]; found {
		s.handle[address] = append(s.handle[address], &router{
			path:    path,
			handler: fn,
		})
	}
}

func (s *Server) Run() {
	if len(s.srv) <= 0 {
		log.Fatal("启动异常，请检查反向代理后端服务设置")
	} else {
		for address, srv := range s.srv {
			log.Println("Starting httpServer at " + address)

			//拿到每一个 http.Server,写入handler
			routers := s.handle[address]
			if len(routers) > 0 {
				for _, v := range routers {
					srv.HandleFunc(v.path, v.handler)
				}
			}

			server := &http.Server{
				Addr:         address,
				WriteTimeout: time.Second * 3,
				Handler:      srv,
			}
			go func() {
				log.Fatal(server.ListenAndServe())
			}()
		}
	}
}

func NewProxy(address string) {
	streamId := "http://127.0.0.1:9091"
	uri, err1 := url.Parse(streamId)
	if err1 != nil {
		log.Println(err1)
	}
	srv := hostProxy(uri)

	log.Fatal(http.ListenAndServe(address, srv.proxy))
}

func hostProxy(target *url.URL) *Server {
	s := &Server{}

	targetQuery := target.RawQuery

	//拿到反向代理的请求参数
	s.director = func(req *http.Request) {
		re, _ := regexp.Compile("^/dir(.*)")
		req.URL.Path = re.ReplaceAllString(req.URL.Path, "$1")
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)

		//拼接Query参数
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "")
		}
	}

	//拿到代理的后端结果
	s.modify = func(res *http.Response) error {

		if res.StatusCode != 200 {
			return errors.New("error statusCode")
		}
		oldPayload, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		//newPayLoad := []byte("hello " + string(oldPayload))
		res.Body = ioutil.NopCloser(bytes.NewBuffer(oldPayload))
		res.ContentLength = int64(len(oldPayload))
		res.Header.Set("Content-Length", fmt.Sprint(len(oldPayload)))
		return nil
	}

	errorHandler := func(res http.ResponseWriter, req *http.Request, err error) {
		res.Write([]byte(err.Error()))
	}
	s.proxy = &httputil.ReverseProxy{
		Director:       s.director,
		ModifyResponse: s.modify,
		ErrorHandler:   errorHandler,
	}
	return s
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
