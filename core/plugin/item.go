package plugin

type PluginItem interface {
	Start() error //开始插件
	Exit()        //退出插件
	Enable() bool //插件开关
	Name() string //服务名称
}
