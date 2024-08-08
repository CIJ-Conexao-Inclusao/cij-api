package config

import "github.com/spf13/viper"

type Config struct {
	DbConnection string `mapstructure:"DSN"`
	SecretKey    string `mapstructure:"SECRET_KEY"`
}

type CloudinaryConfig struct {
	CloudinaryUrl string `mapstructure:"CLOUDINARY_URL"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func LoadCloudinaryConfig(path string) (config CloudinaryConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
