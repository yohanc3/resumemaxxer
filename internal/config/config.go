package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

var Cfg Config

type Config struct {
	DBPassword               string `env:"DB_PASSWORD"`
	DBUsername               string `env:"DB_USERNAME"`
	DBPort                   string `env:"DB_PORT"`
	DBHost                   string `env:"DB_HOST"`
	DBName                   string `env:"DB_NAME"`
	R2BucketName             string `env:"R2_BUCKET_NAME"`
	R2SecretAccessKey        string `env:"R2_SECRET_ACCESS_KEY"`
	R2AccountID              string `env:"R2_ACCOUNT_ID"`
	R2AccessKeyID            string `env:"R2_ACCESS_KEY_ID"`
	DelegatorIntervalSeconds int    `env:"DELEGATOR_INTERVAL_SECONDS"`
}

func LoadConfig() error {

	if err := env.Parse(&Cfg); err != nil {
		return fmt.Errorf("Parsing env vars failed. %w", err)
	}

	return nil
}
