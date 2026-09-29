package cli

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/moomaideng/bogbank/backend/internal/database"
	"github.com/moomaideng/bogbank/backend/internal/httpserver"
	"github.com/moomaideng/bogbank/backend/internal/objectstorage"
	"github.com/moomaideng/bogbank/backend/services/template/internal/config"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	RunE:  runServe,
}

func init() {
	serveCmd.Flags().StringP("config", "c", "", "config file path")
	if err := serveCmd.MarkFlagRequired("config"); err != nil {
		panic(err)
	}
}

func runServe(cmd *cobra.Command, _ []string) error {
	path, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}

	cfg, err := config.Load(path)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var pingers []httpserver.Pinger
	if cfg.Database != nil {
		db, err := database.NewPostgres(ctx, cfg.Database.DSN)
		if err != nil {
			return err
		}
		defer func() {
			if err := db.Close(); err != nil {
				slog.Error("close postgres", "err", err)
			}
		}()
		pingers = append(pingers, db)
		slog.Info("postgres connected")
	}

	if cfg.S3 != nil {
		client, err := objectstorage.New(ctx, objectstorage.Config{
			Endpoint:        cfg.S3.Endpoint,
			Region:          cfg.S3.Region,
			Bucket:          cfg.S3.Bucket,
			AccessKeyID:     cfg.S3.AccessKeyID,
			SecretAccessKey: cfg.S3.SecretAccessKey,
		})
		if err != nil {
			return err
		}
		pingers = append(pingers, client)
		slog.Info("s3 client ready", "bucket", cfg.S3.Bucket)
	}

	server := httpserver.New(cfg.HTTP.Address, "Bogbank Template", pingers)
	return server.Start(ctx)
}
