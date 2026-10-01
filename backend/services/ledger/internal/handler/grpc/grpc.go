package grpc

import (
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/moomaideng/bogbank/backend/services/ledger/internal/deps"
	ledgerv1 "github.com/moomaideng/bogbank/backend/services/ledger/internal/proto/ledger/v1"
	"github.com/moomaideng/bogbank/backend/services/ledger/internal/usecase"
)

type Server struct {
	ledgerv1.UnimplementedLedgerServiceServer
	categories *usecase.CategoryUsecase
}

func Register(server *grpc.Server, deps deps.Deps) {
	ledgerv1.RegisterLedgerServiceServer(server, &Server{
		categories: deps.Categories,
	})
}

func parseUserID(userIDStr string) (uuid.UUID, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil || userID == uuid.Nil {
		return uuid.Nil, status.Error(codes.InvalidArgument, "user id is invalid")
	}
	return userID, nil
}

func parseCategoryID(categoryIDStr string) (uuid.UUID, error) {
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil || categoryID == uuid.Nil {
		return uuid.Nil, status.Error(codes.InvalidArgument, "category id is invalid")
	}
	return categoryID, nil
}

func statusErr(err error) error {
	switch {
	case errors.Is(err, usecase.ErrNotFound):
		return status.Error(codes.NotFound, "category not found")
	case errors.Is(err, usecase.ErrConflict):
		return status.Error(codes.AlreadyExists, "category name already exists")
	case errors.Is(err, usecase.ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, "invalid argument")
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
