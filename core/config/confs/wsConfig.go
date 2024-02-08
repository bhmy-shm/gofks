package confs

type (
	WsConfig struct {
		WS *ws `yaml:"ws"`
	}
	ws struct {
		RpcRouter bool   `yaml:"rpcRouter"`
		SendBytes uint32 `yaml:"sendBytes"`
		NodeId    int64  `yaml:"nodeId"`
		MaxConn   int64  `yaml:"maxConn"`
	}
)

func (c *WsConfig) IsLoad() bool {
	if c.WS == nil {
		return false
	}
	return c.WS.RpcRouter
}
