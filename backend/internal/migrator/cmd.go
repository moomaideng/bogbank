// Package migrator is the shared goose cobra command.
package migrator

import (
	"context"
	"database/sql"
	"io/fs"
	"strings"
	"time"

	"github.com/moomaideng/bogbank/internal/database"
	"github.com/spf13/cobra"
)

// OpenDB returns a database connection. The caller owns flags and config.
// The returned function closes the connection.
type OpenDB func(ctx context.Context, cmd *cobra.Command) (*sql.DB, func(), error)

// NewCommand returns migrate up|down|reset|create. migrations is the embedded
// SQL set. dir is the on-disk directory create writes into (repo-root relative).
func NewCommand(migrations fs.FS, dir string, openDB OpenDB) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate the database",
	}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "up",
			Short: "Apply pending migrations",
			RunE: func(cmd *cobra.Command, _ []string) error {
				return withDB(cmd, openDB, func(ctx context.Context, db *sql.DB) error {
					return database.Up(ctx, db, migrations)
				})
			},
		},
		&cobra.Command{
			Use:   "down",
			Short: "Roll back the last migration",
			RunE: func(cmd *cobra.Command, _ []string) error {
				return withDB(cmd, openDB, func(ctx context.Context, db *sql.DB) error {
					return database.Down(ctx, db, migrations)
				})
			},
		},
		&cobra.Command{
			Use:   "reset",
			Short: "Roll back every migration",
			RunE: func(cmd *cobra.Command, _ []string) error {
				return withDB(cmd, openDB, func(ctx context.Context, db *sql.DB) error {
					err := database.Reset(ctx, db, migrations)
					// A database that has never been migrated has no goose version table.
					if err != nil && strings.Contains(err.Error(), "failed to get status of migrations") {
						return nil
					}
					return err
				})
			},
		},
		&cobra.Command{
			Use:   "create NAME",
			Short: "Write a new SQL migration file",
			Args:  cobra.ExactArgs(1),
			RunE: func(_ *cobra.Command, args []string) error {
				return database.Create(dir, args[0])
			},
		},
	)
	return cmd
}

func withDB(cmd *cobra.Command, openDB OpenDB, fn func(ctx context.Context, db *sql.DB) error) error {
	ctx, cancel := context.WithTimeout(cmd.Context(), 60*time.Second)
	defer cancel()

	db, closeDB, err := openDB(ctx, cmd)
	if err != nil {
		return err
	}
	defer closeDB()
	return fn(ctx, db)
}
