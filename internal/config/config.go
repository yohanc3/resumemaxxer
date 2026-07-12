package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

var Cfg Config

type Config struct {
	DBPassword string `env:"DB_PASSWORD,required"`
	DBUsername string `env:"DB_USERNAME,required"`
	DBPort     string `env:"DB_DOCKER_PORT,required"`
	DBDriver   string `env:"DB_DRIVER,required"`
}

func LoadConfig() error {

	if err := env.Parse(&Cfg); err != nil {
		return fmt.Errorf("Parsing env vars failed. %w", err)
	}

	return nil
}
