package pkg

import (
	"log"
	"testing"
)

//读取文件
func TestLoadFile(t *testing.T) {
	f, err := LoadFile()
	if err != nil {
		log.Fatalln(err)
	}
	f.YamlMerge()

	//开启监听
	go ReadWatcher(f)

	d, err := GetPath("Server", "port").Int()
	if err != nil {
		log.Println(err)
	}
	log.Println("data=", d)
	select {}
}
