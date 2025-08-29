package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"sync"
	"time"
)

var (
	once     sync.Once
	cfg      *Config
	cfgErr   error
	mu       sync.RWMutex
	watchers []chan *Config
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

func (c *Config) Watch() <-chan *Config {
	ch := make(chan *Config)
	mu.Lock()
	watchers = append(watchers, ch)
	mu.Unlock()
	return ch
}

func Get(path string) (*Config, error) {
	once.Do(func() {
		cfg, cfgErr = loadConfig(path)
	})
	return cfg, cfgErr
}

func loadConfig(path string) (*Config, error) {
	log.Println("loading config")

	v, err := setupViper(path)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	if err := v.Unmarshal(&c); err != nil {
		return nil, err
	}

	v.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("config file changed: %s", e.Name)

		newCfg := &Config{}
		if err := v.Unmarshal(&newCfg); err != nil {
			log.Printf("failed to reload config: %v", err)
			return
		}

		mu.Lock()
		cfg = newCfg
		for _, w := range watchers {
			select {
			case w <- newCfg:
			default:
			}
		}
		mu.Unlock()
	})

	return c, nil
}

func setupViper(path string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigFile(path)

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	v.WatchConfig()

	return v, nil
}
