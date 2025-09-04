package configs

type (
	Config struct {
		Service  Service  `mapstructure:"service"`
		Database Database `mapstructure:"database"`
		JWT      JWT      `mapstructure:"jwt"`
		Upload   Upload   `mapstructure:"upload"`
	}

	Service struct {
		Port string `mapstructure:"port"`
	}

	Database struct {
		DataSourceName string `mapstructure:"dataSourceName"`
	}

	JWT struct {
		SecretKey string `mapstructure:"secretKey"`
	}

	Upload struct {
		Path string `mapstructure:"path"`
	}
)
