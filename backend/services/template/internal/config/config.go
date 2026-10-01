package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	sharedconfig "github.com/moomaideng/bogbank/backend/internal/config"
)

const envPrefix = "TEMPLATE"

type Config struct {
	HTTP     HTTP      `mapstructure:"http"`
	Database *Database `mapstructure:"database"`
	S3       *S3       `mapstructure:"s3"`
}

type HTTP struct {
	Address string `mapstructure:"address" validate:"required"`
}

type Database struct {
	DSN string `mapstructure:"dsn" validate:"required"`
}

type S3 struct {
	Endpoint        string `mapstructure:"endpoint"          validate:"required"`
	Region          string `mapstructure:"region"`
	Bucket          string `mapstructure:"bucket"            validate:"required"`
	AccessKeyID     string `mapstructure:"access_key_id"     validate:"required"`
	SecretAccessKey string `mapstructure:"secret_access_key" validate:"required"`
}

func Load(path string) (Config, error) {
	cfg, err := sharedconfig.Load[Config](path, envPrefix)
	if err != nil {
		return Config{}, err
	}
	if err := validator.New().Struct(&cfg); err != nil {
		return Config{}, fmt.Errorf("invalid config: %w", err)
	}
	return cfg, nil
}
