package svc

import gofkConf "github.com/bhmy-shm/gofks/core/config"

type CascadeContext struct {
	Config *gofkConf.Config
}

func NewCascadeContext(conf *gofkConf.Config) *CascadeContext {
	return &CascadeContext{
		Config: conf,
	}
}
