package main

import (
	"github.com/bhmy-shm/gofks"
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	"github.com/bhmy-shm/gofks/example/api/controls/account"
	"github.com/bhmy-shm/gofks/example/api/wire"
)

func main() {

	conf := gofkConf.Load()

	gofks.Ignite("/v1").
		LoadWatch(conf).
		WireApply(
			wire.NewServiceWrite(conf), //自实现依赖注入
		).
		Mount(account.UserController()).
		Launch()
}
