package pkg

import (
	"github.com/bhmy-shm/gofks/core/config/confs"
	"log"
	"testing"
)

const confPath = "./application.yaml"

func TestConfigLoad(t *testing.T) {

	conf1 := Load(WithPath(confPath))
	log.Println("config1", *conf1.LogConfig, *conf1.ServerConfig)

	conf2 := &confs.LogConfig{}
	LoadConf(conf2, WithPath(confPath))

	log.Println("config2", *conf2)
}

func TestConfigLoadFile(t *testing.T) {

}
