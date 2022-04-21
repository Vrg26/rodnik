package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port        string `yaml:"port" env:"HTTP_PORT" envDefault:"8080"`
	LogLevel    string `yaml:"log_level" env:"LOG_LEVEL" envDefault:"info"`
	DataBaseDSN string `yaml:"data_base_dsn" env:"DATABASE_DSN"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
