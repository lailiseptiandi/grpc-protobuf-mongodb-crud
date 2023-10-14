package config

import "github.com/spf13/viper"

type Config struct {
	DBUri  string `mapstructure:"MONGODB_LOCAL_URI"`
	Port   string `mapstructure:"PORT"`
	DBNAME string `mapstructure:"MONGODB_NAME"`

	GrpcServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`

	Origin string `mapstructure:"CLIENT_ORIGIN"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return config, nil

}
