package pkg

import (
	"log"
	"testing"
)

func TestWrapValue(t *testing.T) {

	ss := newValue("111").Value().Unwrap()
	log.Printf("end:%v, %T\n", ss, ss)

	ii := newValue(2134).Value().Unwrap()
	log.Printf("end:%v, %T\n", ii, ii)

	ff := newValue(2.1).Value().Unwrap()
	log.Printf("end:%v, %T\n", ff, ff)
}

func TestGetPath(t *testing.T) {

	c := Load()
	log.Println(c)

	ss := GetPath[string]("server", "name").Value().Unwrap()
	log.Println("getPath server-name:", ss)

	ss2 := GetPath[int]("server", "name").Value().Unwrap()
	log.Println("getPath server-name 2:", ss2)
}
