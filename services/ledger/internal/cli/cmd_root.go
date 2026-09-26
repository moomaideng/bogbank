package cli

import (
	"github.com/moomaideng/bogbank/internal/migrator"
	"github.com/moomaideng/bogbank/services/ledger/internal/migrations"
	"github.com/spf13/cobra"
)

const (
	configFlag    = "config"
	migrationsDir = "services/ledger/internal/migrations" // relative to repo root
)

var rootCmd = &cobra.Command{
	Use:   "ledger",
	Short: "Bogbank ledger service",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringP(configFlag, "c", "", "config file path")
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(migrator.NewCommand(migrations.Migrations, migrationsDir, openDB))
}
