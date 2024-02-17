package controllers

import (
	"github.com/bhmy-shm/gofks"
	"github.com/bhmy-shm/gofks/example/k8s-manager/wire"
	"github.com/gin-gonic/gin"
)

type DeploymentCtl struct {
	*wire.ServiceWire `inject:"-"`
}

func NewDeploymentCtl() *DeploymentCtl {
	return &DeploymentCtl{}
}

func (d *DeploymentCtl) Build(gofk *gofks.Gofk) {
	deploy := gofk.Group("/deployment")
	deploy.GET("/list", d.GetList)
}

func (d *DeploymentCtl) Name() string {
	return "deploymentCtl"
}

func (d *DeploymentCtl) GetList(c *gin.Context) {
	ns, ok := c.GetQuery("namespace")
	if !ok {
		InternalResp(c, RespField("reason", "not Query namespace"))
		return
	}

	list, err := d.Context().DepSvr.List(ns)
	if err != nil {
		InternalResp(c, RespField(err.Error(), d.Name()))
		return
	}

	SuccessResp(c, list)
}
