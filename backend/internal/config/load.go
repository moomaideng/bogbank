// Package config loads a YAML file and overlays environment variables.
package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Load reads path as YAML into T, then overlays environment variables.
// A config key http.address is replaced by PREFIX_HTTP_ADDRESS when that
// variable is set. Keys absent from the file are ignored.
func Load[T any](path, envPrefix string) (T, error) {
	var cfg T
	if envPrefix == "" {
		return cfg, errors.New("env prefix is required")
	}

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(path)
	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return cfg, fmt.Errorf("read config: %w", err)
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("unmarshal config: %w", err)
	}
	return cfg, nil
}
