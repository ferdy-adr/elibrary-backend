package configs

import (
	"log"

	"github.com/spf13/viper"
)

var config *Config

type option struct {
	configFolders []string
	configFile    string
	configType    string
}

func Init(opts ...Option) error {
	opt := &option{
		configFolders: getDefaultConfigFolder(),
		configFile:    getDefaultConfigFile(),
		configType:    getDefaultConfigType(),
	}

	for _, optFunc := range opts {
		optFunc(opt)
	}

	for _, configFolder := range opt.configFolders {
		viper.AddConfigPath(configFolder)
	}

	viper.SetConfigName(opt.configFile)
	viper.SetConfigType(opt.configType)
	viper.AutomaticEnv()

	// Set environment variable mappings for Railway
	viper.SetEnvPrefix("")
	viper.BindEnv("service.port", "PORT")
	viper.BindEnv("database.dataSourceName", "DATABASE_URL")
	viper.BindEnv("jwt.secretKey", "JWT_SECRET_KEY")
	viper.BindEnv("upload.path", "UPLOAD_PATH")

	config = new(Config)

	err := viper.ReadInConfig()
	if err != nil {
		// If config file not found, try to use environment variables only
		log.Println("Config file not found, using environment variables only")
		return viper.Unmarshal(&config)
	}
	return viper.Unmarshal(&config)
}

type Option func(*option)

func getDefaultConfigFolder() []string {
	return []string{"./internal/configs"}
}

func getDefaultConfigFile() string {
	return "config"
}
func getDefaultConfigType() string {
	return "yaml"
}

func WithConfigFolder(configFolders []string) Option {
	return func(opt *option) {
		opt.configFolders = configFolders
	}
}

func WithConfigFile(configFile string) Option {
	return func(opt *option) {
		opt.configFile = configFile
	}
}

func WithConfigType(configType string) Option {
	return func(opt *option) {
		opt.configType = configType
	}
}

func Get() *Config {
	if config == nil {
		config = &Config{}
	}
	return config
}
