package cli

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/moomaideng/bogbank/backend/internal/database"
	"github.com/moomaideng/bogbank/backend/internal/httpserver"
	"github.com/moomaideng/bogbank/backend/services/ledger/internal/deps"
	grpchandler "github.com/moomaideng/bogbank/backend/services/ledger/internal/handler/grpc"
	resthandler "github.com/moomaideng/bogbank/backend/services/ledger/internal/handler/rest"
	"github.com/moomaideng/bogbank/backend/services/ledger/internal/migrations"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP and gRPC servers",
	RunE:  runServe,
}

func runServe(cmd *cobra.Command, _ []string) error {
	cfg, err := loadConfig(cmd)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := database.NewPostgres(ctx, cfg.Database.DSN)
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("close postgres", "err", err)
		}
	}()

	migrateCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	if err := database.Up(migrateCtx, db.DB, migrations.Migrations); err != nil {
		return err
	}
	slog.Info("postgres connected")

	dependencies := deps.New(db)

	lis, err := net.Listen("tcp", cfg.GRPC.Address)
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	grpchandler.Register(grpcServer, dependencies)

	if cfg.GRPC.ReflectionEnabled {
		reflection.Register(grpcServer)
	}

	go func() {
		slog.Info("grpc listening", "address", cfg.GRPC.Address)
		if err := grpcServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			slog.Error("grpc stopped", "err", err)
		}
	}()
	defer grpcServer.GracefulStop()

	server := httpserver.New(cfg.HTTP.Address, "Bogbank Ledger", []httpserver.Pinger{db})
	resthandler.Register(server.HumaAPI(), dependencies)
	return server.Start(ctx)
}
