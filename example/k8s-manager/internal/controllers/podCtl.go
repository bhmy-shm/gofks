package controllers

//type PodCtl struct {
//	PodService *service.PodService `inject:"-"`
//}
//
//func NewPodCtl() *PodCtl {
//	return &PodCtl{}
//}
//
//func (this *PodCtl) GetAll(c *gin.Context) {
//	ns, ok := c.GetQuery("namespace")
//	if !ok {
//		gofk.InternalResp(c, errorx.BadRequest("namespace", "get query is failed"))
//		return
//	}
//	gofk.Successful(c, this.PodService.ListByNs(ns))
//}
//
//func (this *PodCtl) Build(gofk *gofk.Gofk) {
//	pods := gofk.Group("pods")
//	pods.GET("/list", this.GetAll)
//}
//
//func (*PodCtl) Name() string {
//	return "PodCtl"
//}
