package config

import (
	"github.com/bhmy-shm/gofks/pkg/errorx"
	"io/ioutil"
	"log"
	"os"
)

type FileMap map[string]interface{}

var GlobalConf = make(FileMap)

type File struct {
	path    string
	opts    Options
	set     *ChangeSet
	confMap FileMap
}

func LoadFile() (*File, error) {
	dir, _ := os.Getwd()
	fname := dir + "/application.yaml"

	f := File{path: fname}
	f.opts = newOptions()
	f.confMap = make(FileMap)

	return f.Read()
}

func (f *File) String() string {
	return "File"
}

func (f *File) Read() (*File, error) {
	fh, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	//读取文件内容,必须用这个ReadFile
	b, err := ioutil.ReadFile(f.path)
	if err != nil {
		return nil, errorx.FileReadFail
	}
	log.Println("file Read 读取到的文件内容：", string(b))

	//判断文件是否存在，拿到文件的详细信息
	info, err := fh.Stat()
	if err != nil {
		return nil, errorx.FileNotExist
	}

	cs := &ChangeSet{
		Format:    f.path,
		Source:    f.String(),
		Timestamp: info.ModTime(),
		Data:      b,
	}
	f.set = cs
	cs.Checksum = cs.Sum()
	return f, nil
}

func (f *File) Watch() (Watcher, error) {
	if _, err := os.Stat(f.path); err != nil {
		return nil, err
	}
	return newWatcher(f)
}

func (f *File) Reload() *File {
	newf, err := f.Read()
	if err != nil {
		log.Fatalln("重新读取配置文件失败：", err)
	}
	newf.YamlMerge()
	return newf
}

func (f *File) YamlMerge() {
	encode := f.opts.Encoding["yaml"]

	err := encode.Decode(f.set.Data, &f.confMap)
	if err != nil {
		log.Fatalln(err)
	}
	GlobalConf = f.confMap
}

func (f *File) GetConf() FileMap {
	return f.confMap
}


