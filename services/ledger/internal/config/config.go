package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	sharedconfig "github.com/moomaideng/bogbank/internal/config"
)

const envPrefix = "LEDGER"

type Config struct {
	HTTP     HTTP     `mapstructure:"http"`
	GRPC     GRPC     `mapstructure:"grpc"`
	Database Database `mapstructure:"database"`
}

type HTTP struct {
	Address string `mapstructure:"address" validate:"required"`
}

type GRPC struct {
	Address           string `mapstructure:"address" validate:"required"`
	ReflectionEnabled bool   `mapstructure:"reflection_enabled"`
}

type Database struct {
	DSN string `mapstructure:"dsn" validate:"required"`
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
