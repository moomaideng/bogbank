// Package usecase holds ledger rules and the repository interface they call.
package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	CreatedAt time.Time
}

type CategoryRepo interface {
	CreateOne(ctx context.Context, category Category) (Category, error)
	FindOneByUserIDAndCategoryID(ctx context.Context, userID, categoryID uuid.UUID) (Category, error)
	FindManyByUserID(ctx context.Context, userID uuid.UUID) ([]Category, error)
	UpdateNameByUserIDAndCategoryID(ctx context.Context, userID, categoryID uuid.UUID, name string) (Category, error)
	DeleteOneByUserIDAndCategoryID(ctx context.Context, userID, categoryID uuid.UUID) error
}

type CategoryUsecase struct {
	repo CategoryRepo
}

func NewCategoryUsecase(repo CategoryRepo) *CategoryUsecase {
	return &CategoryUsecase{repo: repo}
}

func (s *CategoryUsecase) Create(ctx context.Context, userID uuid.UUID, name string) (Category, error) {
	name, err := validCategoryName(name)
	if err != nil {
		return Category{}, err
	}
	return s.repo.CreateOne(ctx, Category{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      name,
		CreatedAt: time.Now().UTC(),
	})
}

func (s *CategoryUsecase) Get(ctx context.Context, userID, categoryID uuid.UUID) (Category, error) {
	return s.repo.FindOneByUserIDAndCategoryID(ctx, userID, categoryID)
}

func (s *CategoryUsecase) List(ctx context.Context, userID uuid.UUID) ([]Category, error) {
	return s.repo.FindManyByUserID(ctx, userID)
}

func (s *CategoryUsecase) Update(ctx context.Context, userID, categoryID uuid.UUID, name string) (Category, error) {
	name, err := validCategoryName(name)
	if err != nil {
		return Category{}, err
	}
	return s.repo.UpdateNameByUserIDAndCategoryID(ctx, userID, categoryID, name)
}

func (s *CategoryUsecase) Delete(ctx context.Context, userID, categoryID uuid.UUID) error {
	return s.repo.DeleteOneByUserIDAndCategoryID(ctx, userID, categoryID)
}

func validCategoryName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", fmt.Errorf("name is required: %w", ErrInvalidArgument)
	}
	return name, nil
}
