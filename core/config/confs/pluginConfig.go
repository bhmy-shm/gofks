package confs

// PluginConfig
type (
	PluginConfig struct {
		Plugin plugin `yaml:"plugin"`
	}
	plugin struct {
		ServiceName     string `yaml:"serviceName"`
		NodeId          int    `yaml:"nodeId"`
		MonitorInterval int    `yaml:"monitorInterval"`
		MonitorEnable   bool   `yaml:"monitorEnable"`
		RegisterEnable  bool   `yaml:"registerEnable"`
	}
)

func (c *PluginConfig) IsLoad() bool {
	return true
}

func (c *PluginConfig) MonitorInterval() int {
	return c.Plugin.MonitorInterval
}
