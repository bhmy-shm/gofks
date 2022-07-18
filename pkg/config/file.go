package config

import (
	"io/ioutil"
	"log"
	"os"
)

type fileMap map[string]interface{}

type file struct {
	path string
	code Encoder
	set  *ChangeSet
}

func LoadFile() (*file, error) {
	dir, _ := os.Getwd()
	fname := dir + "/application.yaml"

	f := file{path: fname}
	return f.read()
}

func (f *file) String() string {
	return "file"
}

func (f *file) read() (*file, error) {
	fh, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	b, err := ioutil.ReadAll(fh)
	if err != nil {
		return nil, err
	}
	info, err := fh.Stat()
	if err != nil {
		return nil, err
	}

	cs := &ChangeSet{
		Format:    f.path,
		Source:    f.String(),
		Timestamp: info.ModTime(),
		Data:      b,
	}
	f.set = cs
	//cs.Checksum = cs.Sum()
	return f, nil
}

func (f *file) YamlMerge() fileMap {
	f.code = NewYamlEncoder()
	fm := make(fileMap)

	err := f.code.Decode(f.set.Data, &fm)
	if err != nil {
		log.Fatalln(err)
	}
	return fm
}
