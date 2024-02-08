package confs

// RpcServerConf is a rpc server config.
type (
	RpcServerConf struct {
		RpcServer rpcServer `yaml:"rpcServer"`
	}
	rpcServer struct {
		Server        *ServerConfig
		StrictControl bool  `yaml:"strictControl"`
		Timeout       int64 `yaml:"timeout"`      //default 2000
		CpuThreshold  int64 `yaml:"cpuThreshold"` //default=900,range=[0:1000]
	}
)

func (c *RpcServerConf) IsLoad() bool {
	return true
}

// RpcClientConf is a rpc client config.
type (
	RpcClientConf struct {
		RpcClient *rpcClient `yaml:"rpcClient"`
	}
	rpcClient struct {
		Endpoints []string `yaml:"endpoints"`
		Target    string   `yaml:"target"`
		App       string   `yaml:"app"`
		Token     string   `yaml:"token"`
		NonBlock  bool     `yaml:"nonBlock"`
		Timeout   int64    `yaml:"timeout"` //default=2000
	}
)

func (c *RpcClientConf) IsLoad() bool {
	if c.RpcClient == nil {
		return false
	}
	return len(c.RpcClient.Target) > 0
}

func (c *RpcClientConf) HasCredential() bool {
	return len(c.RpcClient.App) > 0 && len(c.RpcClient.Token) > 0
}

// BuildTarget builds the rpc target from the given config.
func (c *RpcClientConf) BuildTarget() (string, error) {

	return c.RpcClient.Target, nil
	//if len(cc.Endpoints) > 0 {
	//	return resolver.BuildDirectTarget(cc.Endpoints), nil
	//} else if len(cc.Target) > 0 {
	//	return cc.Target, nil
	//}
	//
	//if err := cc.Etcd.Validate(); err != nil {
	//	return "", err
	//}
	//
	//if cc.Etcd.HasAccount() {
	//	discov.RegisterAccount(cc.Etcd.Hosts, cc.Etcd.User, cc.Etcd.Pass)
	//}
	//if cc.Etcd.HasTLS() {
	//	if err := discov.RegisterTLS(cc.Etcd.Hosts, cc.Etcd.CertFile, cc.Etcd.CertKeyFile,
	//		cc.Etcd.CACertFile, cc.Etcd.InsecureSkipVerify); err != nil {
	//		return "", err
	//	}
	//}

	//return resolver.BuildDiscovTarget(cc.Etcd.Hosts, cc.Etcd.Key), nil
}

func (c *RpcClientConf) NonBlock() bool {
	return c.RpcClient.NonBlock
}

func (c *RpcClientConf) Timeout() int64 {
	return c.RpcClient.Timeout
}
