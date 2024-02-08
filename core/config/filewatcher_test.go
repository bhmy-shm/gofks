package pkg

import (
	"log"
	"testing"
)

func ReadWatcher2(f *File, conf *Config) {
	w, _ := f.Watch()
	newf, err := w.Next(conf)

	if err != nil {
		log.Fatalln("watch faile =", err)
	}
	f = newf
	if f.confMap != nil {
		log.Println("watch filed after", f.GetString())
		ReadWatcher2(f, conf)
	} else {
		log.Fatalln("监听后的内容出现异常")
	}
}

func TestReadWatcher(t *testing.T) {

	conf := Load()

	//读取并加载当前配置文件
	f, err := LoadFile()
	if err != nil {
		log.Println(err)
	}

	log.Println("before conf:", f.GetString())

	//对当前配置文件开启监听功能
	ReadWatcher2(f, conf)
}
