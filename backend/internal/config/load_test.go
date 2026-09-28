package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnvOverridesYAML(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	body := []byte("http:\n  address: \":8080\"\ndatabase:\n  dsn: \"hello-world\"\n")
	if err := os.WriteFile(path, body, 0o644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("TEMPLATE_HTTP_ADDRESS", ":9090")

	type httpCfg struct {
		Address string `mapstructure:"address"`
	}
	type dbCfg struct {
		DSN string `mapstructure:"dsn"`
	}
	type cfg struct {
		HTTP     httpCfg `mapstructure:"http"`
		Database dbCfg   `mapstructure:"database"`
	}

	got, err := Load[cfg](path, "TEMPLATE")
	if err != nil {
		t.Fatal(err)
	}
	if got.HTTP.Address != ":9090" {
		t.Fatalf("address %q", got.HTTP.Address)
	}
	if got.Database.DSN != "hello-world" {
		t.Fatalf("dsn %q", got.Database.DSN)
	}
}

func TestLoadRequiresPrefix(t *testing.T) {
	if _, err := Load[struct{}]("unused.yaml", ""); err == nil {
		t.Fatal("expected error")
	}
}
