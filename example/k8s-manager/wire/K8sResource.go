package wire

type RESTMapper interface {
	TestInter()
}

type restM struct {
}

func (*restM) TestInter() {}

type K8sRestMapper struct {
	//Client *kubernetes.Clientset `inject:"-"`
}

func NewK8sResource() *K8sRestMapper {
	return &K8sRestMapper{}
}

// RestMapper 所有api groupResource
// 用于初始化和返回一个 RESTMapper，用于将 API 资源映射到 REST 路径，以便在后续的操作中能够根据 API 资源信息进行 REST 请求。
// 通过获取 API 组资源信息并创建 RESTMapper，可以帮助在 Kubernetes 集群中进行资源操作时更方便地进行资源路径映射和管理
func (res *K8sRestMapper) RestMapper() RESTMapper {
	return &restM{}
}
