package plugin

import (
	"context"
	"github.com/bhmy-shm/gofks/core/register"
	"os"
	"time"
)

type OptionFunc func(o *pluginOptions)

type pluginOptions struct {
	id       string
	name     string
	version  string
	metadata map[string]string

	monitorEnable bool //监控开关
	//todo 监控事件驱动

	registerEnable   bool                      //服务注册开关
	registrar        register.Registrar        //服务注册
	instance         *register.ServiceInstance //服务信息
	registrarTimeout time.Duration             //服务注册刷新时间

	ctx         context.Context
	sigs        []os.Signal
	stopTimeout time.Duration
}

// ID with service id.
func ID(id string) OptionFunc {
	return func(o *pluginOptions) { o.id = id }
}

// Name with service name.
func Name(name string) OptionFunc {
	return func(o *pluginOptions) { o.name = name }
}

// Version with service version.
func Version(version string) OptionFunc {
	return func(o *pluginOptions) { o.version = version }
}

// Metadata with service metadata.
func Metadata(md map[string]string) OptionFunc {
	return func(o *pluginOptions) { o.metadata = md }
}

// Context with service context.
func Context(ctx context.Context) OptionFunc {
	return func(o *pluginOptions) { o.ctx = ctx }
}

// Signal with exit signals.
func Signal(sigs ...os.Signal) OptionFunc {
	return func(o *pluginOptions) { o.sigs = sigs }
}

// Registrars with service registry.
func Registrars(r register.Registrar) OptionFunc {
	return func(o *pluginOptions) { o.registrar = r }
}

// RegistrarTimeout with registrar timeout.
func RegistrarTimeout(t time.Duration) OptionFunc {
	return func(o *pluginOptions) { o.registrarTimeout = t }
}

// TODO StopTimeout with app stop timeout.
//func StopTimeout(t time.Duration) OptionFunc {
//	return func(o *pluginOptions) { o.stopTimeout = t }
//}

// TODO Monitor未实现
func Monitor(enable bool) OptionFunc {
	return func(o *pluginOptions) {
		o.monitorEnable = enable
	}
}
