package confs

// MqConfig
type (
	MqConfig struct {
		Rabbit rabbitMq `yaml:"rabbitMq"`
		Redis  redisMq  `yaml:"redisMq"`
	}
	rabbitMq struct {
		Address   string   `yaml:"address"`
		Password  string   `yaml:"password"`
		QueueName string   `yaml:"queueName"`
		Exchange  string   `yaml:"exchange"`
		Key       []string `yaml:"keys"`
	}
	redisMq struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
	}
)

func (c *MqConfig) IsLoad() bool {
	return true
}
