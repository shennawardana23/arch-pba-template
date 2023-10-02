package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type (
	Config struct {
		Server   ServerConfig
		Database DatabaseConfig
		Log      LogConfig
	}

	ServerConfig struct {
		Host string
		Port int
	}

	DatabaseConfig struct {
		Host     string
		Port     int
		DbName   string
		Username string
		Password string
	}

	LogConfig struct {
		Level string
	}
)

var config Config

func Init(env string) {

	if len(env) == 0 {
		env = "local"
	}

	viper.AddConfigPath("config/")
	viper.SetConfigName("config-" + env)
	viper.SetConfigType("yaml")

	// get application config
	err := viper.ReadInConfig()
	panicOnError(err)

	err = viper.Unmarshal(&config)
	panicOnError(err)

	err = viper.Unmarshal(&config)
	panicOnError(err)

	viper.OnConfigChange(func(e fsnotify.Event) {
		viper.Unmarshal(&config)
	})
	viper.WatchConfig()
}

func Get() Config {
	return config
}

func panicOnError(err error) {
	if err != nil {
		log.Printf("panic on config %v", err)
		panic(err)
	}
}
