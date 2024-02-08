package pkg

import (
	"fmt"
	"github.com/bhmy-shm/gofks/core/errorx"
	"github.com/bhmy-shm/gofks/core/utils/hash"
	"log"
	"os"
	"time"
)

const (
	defaultConfName = "/application.yaml"
)

var GlobalConf = make(FileMap)

type (
	Options struct {
		path     string
		yaml     []byte
		Encoding map[string]Encoder
	}
	OptionFunc func(o *Options)

	FileMap map[string]interface{}

	File struct {
		opts    *Options
		set     *ChangeSet
		confMap FileMap
	}
)

func defaultFile() *File {
	dir, _ := os.Getwd()
	filePath := dir + defaultConfName

	fmt.Println("loadFile path:", dir)
	return &File{
		opts: &Options{
			Encoding: map[string]Encoder{
				"json": NewJsonEncoder(),
				"yaml": NewYamlEncoder(),
			},
			path: filePath,
		},
		confMap: make(FileMap),
	}
}

func LoadFile(opts ...OptionFunc) (*File, error) {

	f := defaultFile()

	for _, o := range opts {
		o(f.opts)
	}

	//读取文件
	return f.Read()
}

func WithPath(path string) OptionFunc {
	return func(o *Options) {
		o.path = path
	}
}

func (f *File) Read() (*File, error) {
	fh, err := os.Open(f.opts.path)
	if err != nil {
		errorx.Fatal(err, "读取配置文件失败")
	}
	defer fh.Close()

	//读取文件内容,必须用这个ReadFile
	b, err := os.ReadFile(f.opts.path)
	if err != nil {
		return nil, errorx.Wrap(err, "[config=file] readFile failed：")
	}
	f.opts.yaml = b

	//判断文件是否存在，拿到文件的详细信息
	info, err := fh.Stat()
	if err != nil {
		return nil, errorx.Wrap(err, "[config-file] fn.Stat failed:%v", err)
	}

	cs := &ChangeSet{
		Format:    f.opts.path,
		Source:    f.GetString(),
		Timestamp: info.ModTime(),
		Data:      b,
	}
	f.set = cs
	cs.CheckMd5 = cs.HashMd5()
	return f, nil
}

func (f *File) Watch() (WatcherInter, error) {
	if _, err := os.Stat(f.opts.path); err != nil {
		return nil, err
	}
	return newWatcher(f)
}

func (f *File) Reload(conf *Config) *File {
	newFile, err := f.Read()
	if err != nil {
		log.Fatalln("重新读取配置文件失败：", err)
	}
	newFile.YamlMerge(conf)
	return newFile
}

func (f *File) YamlMerge(conf *Config) bool {
	encode := f.opts.Encoding["yaml"]

	//映射到map中
	err := encode.Decode(f.set.Data, &f.confMap)
	if err != nil {
		errorx.Fatal(err, "yamlConf Merge to confMap failed")
	}

	//映射到配置文件中
	err = conf.loadAll(f)
	if err != nil {
		errorx.Fatal(err, "yamlConf Merge to config failed")
	}

	GlobalConf = f.confMap
	return true
}

func (f *File) GetConf() FileMap {
	return f.confMap
}

func (f *File) GetString() string {
	if f.opts.yaml == nil {
		return ""
	}
	return string(f.opts.yaml)
}

func (f *File) GetBytes() []byte {
	if f.opts.yaml == nil {
		return nil
	}
	return f.opts.yaml
}

// ChangeSet 文件源，data存放文件源所包含的内容，其余字段声明文件的路径等fileInfo信息
type ChangeSet struct {
	Data      []byte
	CheckMd5  string //md5校验值确保唯一
	Format    string
	Source    string
	Timestamp time.Time
}

func (c *ChangeSet) HashMd5() string {
	return hash.Md5Hex(c.Data)
}
