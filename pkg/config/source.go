package config

import (
	"crypto/md5"
	"fmt"
	"time"
)

//读取文件的源内容
type Source interface {
	Read() (*File, error)
	Watch() (Watcher, error)
	String() string

	//Write(*ChangeSet) error
}

type ChangeSet struct {
	Data      []byte
	Checksum  string //md5校验值确保唯一
	Format    string
	Source    string
	Timestamp time.Time
}

func (c *ChangeSet) Sum() string {
	h := md5.New()
	h.Write(c.Data)
	return fmt.Sprintf("%x", h.Sum(nil))
}
