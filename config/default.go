package config

import "github.com/spf13/viper"

type Config struct {
	Port              string `mapstructure:"PORT"`
	GrpcServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`
	Origin            string `mapstructure:"CLIENT_ORIGIN"`

	DBUri      string `mapstructure:"MONGODB_LOCAL_URI"`
	DBNAME     string `mapstructure:"MONGODB_NAME"`
	DBUSERNAME string `mapstructure:"MONGO_INITDB_ROOT_USERNAME"`
	DBPASSWORD string `mapstructure:"MONGO_INITDB_ROOT_PASSWORD"`

	RedisUri string `mapstructure:"REDIS_URL"`
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
