package cli

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/moomaideng/bogbank/backend/internal/database"
	"github.com/moomaideng/bogbank/backend/services/ledger/internal/config"
	"github.com/spf13/cobra"
)

func loadConfig(cmd *cobra.Command) (config.Config, error) {
	path, err := cmd.Flags().GetString(configFlag)
	if err != nil {
		return config.Config{}, err
	}
	if path == "" {
		return config.Config{}, fmt.Errorf("--%s is required", configFlag)
	}
	return config.Load(path)
}

func openDB(ctx context.Context, cmd *cobra.Command) (*sql.DB, func(), error) {
	cfg, err := loadConfig(cmd)
	if err != nil {
		return nil, nil, err
	}
	db, err := database.NewPostgres(ctx, cfg.Database.DSN)
	if err != nil {
		return nil, nil, err
	}
	return db.DB, func() { db.Close() }, nil
}
