package rest

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/moomaideng/bogbank/backend/services/ledger/internal/deps"
	"github.com/moomaideng/bogbank/backend/services/ledger/internal/usecase"
)

func Register(api huma.API, deps deps.Deps) {
	registerCategories(api, deps.Categories)
}

type userInput struct {
	UserID string `path:"user_id"`
}

type idInput struct {
	UserID string `path:"user_id"`
	ID     string `path:"id"`
}

func parseUserID(userIDStr string) (uuid.UUID, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil || userID == uuid.Nil {
		return uuid.Nil, huma.Error400BadRequest("user id is invalid")
	}
	return userID, nil
}

func parseCategoryID(categoryIDStr string) (uuid.UUID, error) {
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil || categoryID == uuid.Nil {
		return uuid.Nil, huma.Error400BadRequest("category id is invalid")
	}
	return categoryID, nil
}

func writeErr(err error) error {
	switch {
	case errors.Is(err, usecase.ErrNotFound):
		return huma.Error404NotFound("category not found")
	case errors.Is(err, usecase.ErrConflict):
		return huma.Error409Conflict("category name already exists")
	case errors.Is(err, usecase.ErrInvalidArgument):
		return huma.Error400BadRequest("invalid argument")
	default:
		return huma.Error500InternalServerError("internal error")
	}
}
