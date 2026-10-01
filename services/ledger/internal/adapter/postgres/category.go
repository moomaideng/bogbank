package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/moomaideng/bogbank/internal/baserepo"
	"github.com/moomaideng/bogbank/services/ledger/internal/usecase"
	"github.com/uptrace/bun"
)

type categoryModel struct {
	bun.BaseModel `bun:"table:categories"`
	ID            uuid.UUID `bun:"id,pk,type:uuid"`
	UserID        uuid.UUID `bun:"user_id,type:uuid,notnull"`
	Name          string    `bun:"name,notnull"`
	CreatedAt     time.Time `bun:"created_at,notnull"`
}

func toCategoryModel(category usecase.Category) categoryModel {
	return categoryModel{
		ID:        category.ID,
		UserID:    category.UserID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
	}
}

func fromCategoryModel(model categoryModel) usecase.Category {
	return usecase.Category{
		ID:        model.ID,
		UserID:    model.UserID,
		Name:      model.Name,
		CreatedAt: model.CreatedAt.UTC(),
	}
}

type CategoryRepo struct {
	base baserepo.BaseRepo[categoryModel]
	exec baserepo.Executor
}

var _ usecase.CategoryRepo = (*CategoryRepo)(nil)

func NewCategoryRepo(db *bun.DB) *CategoryRepo {
	return &CategoryRepo{
		base: baserepo.NewBaseRepo[categoryModel](db),
		exec: baserepo.NewExecutor(db),
	}
}

func (r *CategoryRepo) CreateOne(ctx context.Context, category usecase.Category) (usecase.Category, error) {
	model := toCategoryModel(category)
	if err := r.base.Create(ctx, &model); err != nil {
		return usecase.Category{}, mapWriteErr(err)
	}
	return fromCategoryModel(model), nil
}

func (r *CategoryRepo) FindOneByUserIDAndCategoryID(ctx context.Context, userID, categoryID uuid.UUID) (usecase.Category, error) {
	var model categoryModel
	err := r.exec.Run(ctx, func(idb bun.IDB) error {
		return idb.NewSelect().Model(&model).
			Where("id = ?", categoryID).
			Where("user_id = ?", userID).
			Scan(ctx)
	})
	if err != nil {
		return usecase.Category{}, mapReadErr(err)
	}
	return fromCategoryModel(model), nil
}

func (r *CategoryRepo) FindManyByUserID(ctx context.Context, userID uuid.UUID) ([]usecase.Category, error) {
	var models []categoryModel
	err := r.exec.Run(ctx, func(idb bun.IDB) error {
		return idb.NewSelect().Model(&models).
			Where("user_id = ?", userID).
			Order("created_at ASC", "id ASC").
			Scan(ctx)
	})
	if err != nil {
		return nil, err
	}
	out := make([]usecase.Category, 0, len(models))
	for _, model := range models {
		out = append(out, fromCategoryModel(model))
	}
	return out, nil
}

func (r *CategoryRepo) UpdateNameByUserIDAndCategoryID(ctx context.Context, userID, categoryID uuid.UUID, name string) (usecase.Category, error) {
	var model categoryModel
	err := r.exec.Run(ctx, func(idb bun.IDB) error {
		res, err := idb.NewUpdate().Model(&model).
			Set("name = ?", name).
			Where("id = ?", categoryID).
			Where("user_id = ?", userID).
			Returning("*").
			Exec(ctx)
		if err != nil {
			return mapWriteErr(err)
		}
		n, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if n == 0 {
			return usecase.ErrNotFound
		}
		return nil
	})
	if err != nil {
		return usecase.Category{}, err
	}
	return fromCategoryModel(model), nil
}

func (r *CategoryRepo) DeleteOneByUserIDAndCategoryID(ctx context.Context, userID, categoryID uuid.UUID) error {
	return r.exec.Run(ctx, func(idb bun.IDB) error {
		res, err := idb.NewDelete().Model((*categoryModel)(nil)).
			Where("id = ?", categoryID).
			Where("user_id = ?", userID).
			Exec(ctx)
		if err != nil {
			return err
		}
		n, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if n == 0 {
			return usecase.ErrNotFound
		}
		return nil
	})
}
