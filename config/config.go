package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App   `yaml:"app"`
	Http  `yaml:"http"`
	Log   `yaml:"log"`
	Pg    `yaml:"postgres"`
	Saver `yaml:"saver"`
}

const (
	SAVER_TYPE_GDRIVE = "gdrive"
	SAVER_TYPE_LOCAL  = "local"
)

type (
	App struct {
		Name    string `env-required:"false" yaml:"name"`
		Version string `env-required:"true" yaml:"version"`
	}

	Http struct {
		Port int `env-required:"true" env:"HTTP_PORT" yaml:"port"`
	}

	Log struct {
		Level string `env-required:"true" env:"LOG_LEVEL" yaml:"level"`
	}

	Pg struct {
		URL string `env-required:"true" env:"PG_URL" yaml:"url"`
	}

	Saver struct {
		//Type   string `env-required:"true" env:"SAVER_TYPE" yaml:"type"`
		//APIKey string `env-required:"false" env:"GDRIVE_API_KEY" yaml:"APIKey"`
	}
)

func New(cfgPath string) (*Config, error) {
	cfg := &Config{}
	fmt.Println(cfgPath)
	err := cleanenv.ReadConfig(cfgPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("Couldnt parse coonfig: %v", err)
	}
	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("Couldnt update envs: %v", err)
	}
	return cfg, nil
}
