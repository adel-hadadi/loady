package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"sync"
	"time"
)

var (
	once   sync.Once
	config *Config
)

type Config struct {
	Servers     []string
	Port        int
	Algorithm   string
	Healthcheck struct {
		Duration time.Duration
		Api      string
	}
}

func Get() (*Config, error) {
	var err error

	once.Do(func() {
		log.Println("loading config")
	
		var v *viper.Viper
		v, err = setup()

		err = v.Unmarshal(&config)

		v.OnConfigChange(func(e fsnotify.Event) {
			if err := v.Unmarshal(&config); err != nil {
				log.Fatal(err)
			}
		})
	})

	return config, err
}

func setup() (*viper.Viper, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.WatchConfig()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return viper.GetViper(), nil
}
