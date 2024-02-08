package pkg

import (
	"github.com/bhmy-shm/gofks/core/config/confs"
	"github.com/bhmy-shm/gofks/core/errorx"
	"gopkg.in/yaml.v3"
)

type (
	ConfigItem interface {
		IsLoad() bool
	}
	Config struct {
		*confs.ServerConfig
		*confs.JwtConfig
		*confs.LogConfig
		*confs.RegistryConfig
		*confs.DBConfig
		*confs.RedisConfig
		*confs.MqConfig
		*confs.PluginConfig
		*confs.RpcServerConf
		*confs.RpcClientConf
		*confs.WsConfig
	}
)

func New() *Config {
	return &Config{
		ServerConfig:   new(confs.ServerConfig),
		JwtConfig:      new(confs.JwtConfig),
		LogConfig:      new(confs.LogConfig),
		RegistryConfig: new(confs.RegistryConfig),
		DBConfig:       new(confs.DBConfig),
		RedisConfig:    new(confs.RedisConfig),
		MqConfig:       new(confs.MqConfig),
		PluginConfig:   new(confs.PluginConfig),
		RpcServerConf:  new(confs.RpcServerConf),
		RpcClientConf:  new(confs.RpcClientConf),
		WsConfig:       new(confs.WsConfig),
	}

}

// LoadConf 构建指定的 config配置实例
func LoadConf(config interface{}, opts ...OptionFunc) {
	f := defaultFile()

	for _, o := range opts {
		o(f.opts)
	}

	//读取文件
	f, err := f.Read()
	if err != nil {
		errorx.Fatal(err)
	}

	errorx.Fatal(yaml.Unmarshal(f.GetBytes(), config), "loadConf failed")
}

// Load 构建整个Config
func Load(opts ...OptionFunc) *Config {

	f := defaultFile()

	for _, o := range opts {
		o(f.opts)
	}

	//读取文件
	f, err := f.Read()
	if err != nil {
		errorx.Fatal(err)
	}

	//映射对象实例
	conf := New()

	if err = conf.loadAll(f); err != nil {
		errorx.Fatal(err)
	}
	return conf
}

func (s *Config) GetServer() *confs.ServerConfig {
	return s.ServerConfig
}

func (s *Config) GetLog() *confs.LogConfig {
	return s.LogConfig
}

func (s *Config) GetRegister() *confs.RegistryConfig {
	return s.RegistryConfig
}

func (s *Config) GetAuth() *confs.JwtConfig {
	return s.JwtConfig
}

func (s *Config) GetDB() *confs.DBConfig {
	return s.DBConfig
}

func (s *Config) GetRpcServer() *confs.RpcServerConf {
	s.RpcServerConf.RpcServer.Server = s.ServerConfig
	return s.RpcServerConf
}

func (s *Config) GetRpcClient() *confs.RpcClientConf {
	return s.RpcClientConf
}

func (s *Config) GetRedis() *confs.RedisConfig {
	return s.RedisConfig
}

func (s *Config) GetPlugins() *confs.PluginConfig {
	return s.PluginConfig
}

func (s *Config) GetMq() *confs.MqConfig {
	return s.MqConfig
}

func (s *Config) GetWsCore() *confs.WsConfig {
	return s.WsConfig
}

func (s *Config) loadAll(f *File) error {

	// 定义配置对象的切片
	configs := []interface{}{
		s.ServerConfig,
		s.JwtConfig,
		s.LogConfig,
		s.DBConfig,
		s.MqConfig,
		s.PluginConfig,
		s.RpcServerConf,
		s.RpcClientConf,
		s.RedisConfig,
		s.WsConfig,
	}
	for _, config := range configs {
		err := yaml.Unmarshal(f.GetBytes(), config)
		if err != nil {
			return err
		}
	}
	return nil
}
