package gofk

import (
	"fmt"
	user "github.com/bhmy-shm/gofks/example"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

type HttpServerOne struct {
	addr string
}

func NewHttpServer1() *HttpServerOne {
	return &HttpServerOne{}
}

func (this *HttpServerOne) rootHandle(w http.ResponseWriter, req *http.Request) {
	fmt.Println("hello path=", req.Host)
	upath := fmt.Sprintf("http://%s%s\n", this.addr, req.URL.Path)
	realIP := fmt.Sprintf("RemoteAddr=%s,X-Forwarded-For=%v,X-Real-Ip=%v\n",
		req.RemoteAddr,
		req.Header.Get("X-Forwarded-For"),
		req.Header.Get("X-Real-Ip"))
	io.WriteString(w, upath)
	io.WriteString(w, realIP)
}

func (this *HttpServerOne) UserHandle(w http.ResponseWriter, req *http.Request) {

	upath := fmt.Sprintf("http://%s%s\n", this.addr, req.URL.Path)
	realIP := fmt.Sprintf("RemoteAddr=%s,X-Forwarded-For=%v,X-Real-Ip=%v\n",
		req.RemoteAddr,
		req.Header.Get("X-Forwarded-For"),
		req.Header.Get("X-Real-Ip"))

	log.Println("upath=", upath)
	log.Println("realIp=", realIP)

	uu := &user.UserModel{
		Id:   101,
		Name: "shm",
	}

	io.WriteString(w, uu.Name)
	io.WriteString(w, "用户信息1111111")
}

func (this *HttpServerOne) ErrorHandle(w http.ResponseWriter, req *http.Request) {
	upath := "error handler"
	w.WriteHeader(500)
	io.WriteString(w, upath)
}

func (this *HttpServerOne) SetAddr(server *Server) {
	this.addr = ":9091"
	server.SetServe(this.addr)
}

func (this *HttpServerOne) Build(server *Server) {
	server.HandleFunc(this.addr, "/", this.rootHandle)
	server.HandleFunc(this.addr, "/user", this.UserHandle)
	server.HandleFunc(this.addr, "/error", this.ErrorHandle)
}

type HttpServerTwo struct {
	addr string
}

func NewHttpServer2() *HttpServerTwo {
	return &HttpServerTwo{}
}

func (this *HttpServerTwo) rootHandle(w http.ResponseWriter, req *http.Request) {
	fmt.Println("hello path=", req.Host)
	upath := fmt.Sprintf("http://%s%s\n", this.addr, req.URL.Path)
	realIP := fmt.Sprintf("RemoteAddr=%s,X-Forwarded-For=%v,X-Real-Ip=%v\n",
		req.RemoteAddr,
		req.Header.Get("X-Forwarded-For"),
		req.Header.Get("X-Real-Ip"))
	io.WriteString(w, upath)
	io.WriteString(w, realIP)
}

func (this *HttpServerTwo) UserHandle(w http.ResponseWriter, req *http.Request) {

	upath := fmt.Sprintf("http://%s%s\n", this.addr, req.URL.Path)
	realIP := fmt.Sprintf("RemoteAddr=%s,X-Forwarded-For=%v,X-Real-Ip=%v\n",
		req.RemoteAddr,
		req.Header.Get("X-Forwarded-For"),
		req.Header.Get("X-Real-Ip"))

	log.Println("upath=", upath)
	log.Println("realIp=", realIP)

	uu := &user.UserModel{
		Id:   101,
		Name: "shm",
	}

	io.WriteString(w, uu.Name)
	io.WriteString(w, "用户信息2222222")
}

func (this *HttpServerTwo) ErrorHandle(w http.ResponseWriter, req *http.Request) {
	upath := "error handler"
	w.WriteHeader(500)
	io.WriteString(w, upath)
}

func (this *HttpServerTwo) SetAddr(server *Server) {
	this.addr = ":9092"
	server.SetServe(this.addr)
}

func (this *HttpServerTwo) Build(server *Server) {
	server.HandleFunc(this.addr, "/", this.rootHandle)
	server.HandleFunc(this.addr, "/user", this.UserHandle)
	server.HandleFunc(this.addr, "/error", this.ErrorHandle)
}

func TestWeb(t *testing.T) {

	//开启一个http-server
	NewHttp().Mount(NewHttpServer1(), NewHttpServer2()).Run()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func TestProxy(t *testing.T) {
	//开启代理
	NewProxy("127.0.0.1:2002")

}
