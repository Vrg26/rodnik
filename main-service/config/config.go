package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type (
	Config struct {
		Port      string `yaml:"port" env:"HTTP_PORT" envDefault:"8080"`
		LogLevel  string `yaml:"log_level" env:"LOG_LEVEL" envDefault:"info"`
		PG        `yaml:"postgres"`
		SecretKey string `yaml:"secret_key" env:"SECRET_KEY"`
	}
	PG struct {
		Host     string `yaml:"host" env:"PG_HOST"`
		Port     string `yaml:"port" env:"PG_PORT"`
		User     string `yaml:"user" env:"PG_USER"`
		Password string `yaml:"password" env:"PG_PASSWORD"`
		DBName   string `yaml:"dbname" env:"PG_DB_NAME"`
		SSLMode  string `yaml:"sslmode" env:"PG_SSL_MODE"`
	}
)

var instance *Config
var once sync.Once

func GetConfig() (*Config, error) {
	var err error
	once.Do(func() {
		instance = &Config{}

		err = cleanenv.ReadConfig("./config/config.yml", instance)
		if err != nil {
			return
		}

		err = cleanenv.ReadEnv(instance)
		if err != nil {
			return
		}
	})
	if err != nil {
		return nil, err
	}
	return instance, nil
}
