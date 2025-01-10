package config

import (
	"sync"

	"gitea.com/logicamp/lc"
)

var (
	once    sync.Once
	configs *Config
)

type Config struct {
	PORT        string `env:"PORT" default:"3000"`
	BASE_URL    string `env:"BASE_URL" default:""`
	PG_HOST     string `env:"PG_HOST" default:"localhost:5432"`
	PG_USER     string `env:"PG_USER" default:"postgres"`
	PG_PASSWORD string `env:"PG_PASSWORD" default:"---"`
	PG_DATABASE string `env:"PG_DATABASE" default:"postgres"`
}

func GetConfig() *Config {
	once.Do(func() {
		var err error
		configs, err = lc.GetConfig[Config](&Config{})
		if err != nil {
			panic("environment variable problems")
		}
	})
	return configs
}
