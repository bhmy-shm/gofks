package confs

// DBConfig
type (
	DBConfig struct {
		DB *dbGorm `yaml:"db"`
	}
	dbGorm struct {
		Debug          string `yaml:"debug" description:"Info,Error"`
		Types          string `yaml:"types"`
		DataSourceName string `yaml:"dataSourceName"`
		MaxIdleConns   int    `yaml:"maxIdleConns"`
		MaxOpenConns   int    `yaml:"maxOpenConns"`
		MaxLifetime    int    `yaml:"maxLifetime"`
	}
)

func (c *DBConfig) IsLoad() bool {
	if c.DB == nil {
		return false
	}
	return len(c.DB.DataSourceName) > 0
}
