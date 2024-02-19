package main

import (
	"github.com/bhmy-shm/gofks"
	"github.com/bhmy-shm/gofks/example/k8s-manager/internal/controllers"
	"github.com/bhmy-shm/gofks/example/k8s-manager/internal/middlewares"
	"github.com/bhmy-shm/gofks/example/k8s-manager/wire"
)

func main() {

	gofks.Ignite("/v1", middlewares.OnRequest()).
		WireApply(
			wire.NewK8sResource(),
			wire.NewK8sHandler(),
			//wire.NewK8sConfig(),
			wire.NewK8sMaps(),
			wire.NewServiceWire(),
		).
		Mount(
			controllers.NewDeploymentCtl(),
		).
		Launch()
}
