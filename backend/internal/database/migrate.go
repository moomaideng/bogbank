package database

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
)

// Up applies pending migrations. fsys is the embedded SQL directory; goose
// treats "." as that directory's root.
func Up(ctx context.Context, sqlDB *sql.DB, fsys fs.FS) error {
	if err := bind(fsys); err != nil {
		return err
	}
	if err := goose.UpContext(ctx, sqlDB, "."); err != nil {
		return fmt.Errorf("migrate up: %w", err)
	}
	return nil
}

// Down rolls back one migration.
func Down(ctx context.Context, sqlDB *sql.DB, fsys fs.FS) error {
	if err := bind(fsys); err != nil {
		return err
	}
	if err := goose.DownContext(ctx, sqlDB, "."); err != nil {
		return fmt.Errorf("migrate down: %w", err)
	}
	return nil
}

// Reset rolls back every applied migration.
func Reset(ctx context.Context, sqlDB *sql.DB, fsys fs.FS) error {
	if err := bind(fsys); err != nil {
		return err
	}
	if err := goose.ResetContext(ctx, sqlDB, "."); err != nil {
		return fmt.Errorf("migrate reset: %w", err)
	}
	return nil
}

// Create writes a new SQL migration on disk. dir is the real migrations
// directory. Goose's Create API takes a *sql.DB but does not use it for SQL
// scaffolding, so nil is fine.
func Create(dir, name string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	if err := goose.Create(nil, dir, name, "sql"); err != nil {
		return fmt.Errorf("migrate create: %w", err)
	}
	return nil
}

func bind(fsys fs.FS) error {
	goose.SetBaseFS(fsys)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return nil
}
