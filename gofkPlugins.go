package gofks

import (
	gofkConf "github.com/bhmy-shm/gofks/core/config"
	"github.com/bhmy-shm/gofks/core/plugin"
)

type (
	PluginManager struct {
		plugins []plugin.PluginItem
		opts    []plugin.OptionFunc
		conf    *gofkConf.Config
	}
)

func Plugin(conf *gofkConf.Config) *PluginManager {
	manager := &PluginManager{
		plugins: make([]plugin.PluginItem, 0),
		conf:    conf,
	}

	return manager
}

func (p *PluginManager) Mount(access ...plugin.PluginItem) *PluginManager {
	if len(access) == 0 {
		return p
	}

	for _, v := range access {
		if v.Enable() {
			p.plugins = append(p.plugins, v)
		}
	}
	return p
}

func (p *PluginManager) Attach(opts ...plugin.OptionFunc) *PluginManager {
	p.opts = opts
	return p
}

func (p *PluginManager) Run() {

	server := plugin.NewBase(p.conf.PluginConfig)

	for _, fn := range p.opts {
		fn(server.GetOpts())
	}

	server.Add(p.plugins...)

	server.Run()
}
