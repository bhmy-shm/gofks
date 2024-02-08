package confs

import "time"

// RegistryConfig
type (
	RegistryConfig struct {
		Registry registry `yaml:"registry"`
	}
	registry struct {
		Enable      bool          `yaml:"enable"`
		Namespace   string        `yaml:"namespace"`   //命名空间
		Endpoints   []string      `yaml:"endpoints"`   //连接地址
		DialTimeout time.Duration `yaml:"dialTimeout"` //连接超时时间
		TTL         time.Duration `yaml:"ttl"`         //注册过期时间
		MaxRetry    int           `yaml:"maxRetry"`    //心跳最大重试
	}
)

func (c *RegistryConfig) IsLoad() bool {
	return true
}

func (c *RegistryConfig) IsEnable() bool {
	return c.Registry.Enable
}

func (c *RegistryConfig) MaxRetry() int {
	return c.Registry.MaxRetry
}

func (c *RegistryConfig) Namespace() string {
	return c.Registry.Namespace
}

func (c *RegistryConfig) TTL() time.Duration {
	return c.Registry.TTL
}

func (c *RegistryConfig) DialTimeout() time.Duration {
	return c.Registry.DialTimeout
}

func (c *RegistryConfig) Endpoints() []string {
	return c.Registry.Endpoints
}
