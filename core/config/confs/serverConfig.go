package confs

// ServerConfig
type (
	ServerConfig struct {
		Server server `yaml:"server"`
	}
	server struct {
		Name           string `yaml:"name"`
		Listener       string `yaml:"listener"`
		Mode           string `yaml:"mode"` //debug, test, production
		Timeout        int    `yaml:"timeout"`
		EnableWs       bool   `yaml:"enableWs"`       //true: 启动ws
		EnablePProf    bool   `yaml:"enablePProf"`    //true: 启动pprof
		EnableCron     bool   `yaml:"enableCron"`     //true：启动计划任务
		EnableMetrics  bool   `yaml:"enableMetrics"`  //true：启动rpc的性能检测
		PassEncryption bool   `yaml:"passEncryption"` //true：启动密码加密功能
	}
)

func (c *ServerConfig) IsLoad() bool {
	return true
}

func (c *ServerConfig) Listener() string {
	return c.Server.Listener
}

func (c *ServerConfig) Timeout() int {
	return c.Server.Timeout
}

func (c *ServerConfig) EnableMetrics() bool {
	return c.Server.EnableMetrics
}

func (c *ServerConfig) EnableWs() bool { return c.Server.EnableWs }
