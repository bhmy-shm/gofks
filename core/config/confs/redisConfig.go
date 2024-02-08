package confs

// RedisConfig
type (
	RedisConfig struct {
		Cache *redis `yaml:"cache"`
	}
	redis struct {
		Network   string `yaml:"network"`
		Addr      string `yaml:"addr"`
		Pass      string `yaml:"pass"`
		Type      string `yaml:"type"`
		Tls       bool   `yaml:"tls"`
		MaxIdle   int    `yaml:"maxIdle"`
		MaxActive int    `yaml:"maxActive"`
		Wait      bool   `yaml:"wait"`
	}
)

func (c *RedisConfig) IsLoad() bool {
	if c.Cache == nil {
		return false
	}
	return len(c.Cache.Addr) > 0
}

func (c *RedisConfig) MaxIdle() int {
	return c.Cache.MaxIdle
}

func (c *RedisConfig) MaxActive() int {
	return c.Cache.MaxActive
}

func (c *RedisConfig) Wait() bool {
	return c.Cache.Wait
}

func (c *RedisConfig) Password() string {
	return c.Cache.Pass
}

func (c *RedisConfig) Address() string {
	return c.Cache.Addr
}

func (c *RedisConfig) Network() string {
	return c.Cache.Network
}

func (c *RedisConfig) Type() string {
	return c.Cache.Type
}

func (c *RedisConfig) Tls() bool {
	return c.Cache.Tls
}
