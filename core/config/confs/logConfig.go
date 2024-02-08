package confs

// LogConfig
type (
	LogConfig struct {
		LogC logC `yaml:"log"`
	}
	logC struct {
		ServiceName         string `yaml:"serviceName"`
		Mode                string `yaml:"mode"`
		TimeFormat          string `yaml:"timeFormat"`
		Path                string `yaml:"path"`
		Level               string `yaml:"level"`
		Compress            bool   `yaml:"compress"`
		KeepDays            int    `yaml:"keepDays"`
		StackCoolDownMillis int    `yaml:"stackCoolDownMillis"`
	}
)

func (c *LogConfig) IsLoad() bool {
	return true
}

func (c *LogConfig) Mode() string {
	return c.LogC.Mode
}

func (c *LogConfig) Level() string {
	return c.LogC.Level
}

func (c *LogConfig) Path() string {
	return c.LogC.Path
}

func (c *LogConfig) Compress() bool {
	return c.LogC.Compress
}

func (c *LogConfig) KeepDays() int {
	return c.LogC.KeepDays
}

func (c *LogConfig) StackCoolDownMillis() int {
	return c.LogC.StackCoolDownMillis
}
