package main

import (
	"flag"
	"github.com/bhmy-shm/gofks"
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	"github.com/bhmy-shm/gofks/core/plugin"
	"github.com/bhmy-shm/gofks/core/register"
	"github.com/bhmy-shm/gofks/example/plugins/lib"
)

var confPath string

func init() {
	flag.StringVar(&confPath, "conf", "./application.yaml", "conf")
	flag.Parse()
}

func main() {

	conf := gofkConf.Load(gofkConf.WithPath(confPath))

	gofks.Plugin(conf).
		Attach(
			plugin.ID("001"),
			plugin.Name("FLOW-TEST"),
			plugin.Version("v1"),
			plugin.Registrars(register.NewEtcdRegistry(conf.GetRegister())),
			plugin.Monitor(false),
		).
		Mount(
			lib.NewFlow(),
			lib.NewFlight(),
		).
		Run()

}
