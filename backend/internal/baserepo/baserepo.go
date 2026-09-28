// Package baserepo is the shared CRUD, transaction, and cursor-page API.
package baserepo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var ErrNotFound = errors.New("not found")

// BaseRepo is generic CRUD for a bun model M whose primary key column is id.
type BaseRepo[M any] interface {
	Create(ctx context.Context, model *M) error
	FindByID(ctx context.Context, id uuid.UUID) (*M, error)
	// UpdateByID writes every column on model, including zero values.
	UpdateByID(ctx context.Context, id uuid.UUID, model *M) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type baseRepo[M any] struct {
	executor Executor
}

// NewBaseRepo binds CRUD to db. Queries join an ambient transaction when
// Transactioner.Transaction has stored one on ctx.
func NewBaseRepo[M any](db *bun.DB) BaseRepo[M] {
	return &baseRepo[M]{executor: NewExecutor(db)}
}

func (b *baseRepo[M]) Create(ctx context.Context, model *M) error {
	return b.executor.Run(ctx, func(idb bun.IDB) error {
		_, err := idb.NewInsert().Model(model).Exec(ctx)
		return err
	})
}

func (b *baseRepo[M]) FindByID(ctx context.Context, id uuid.UUID) (*M, error) {
	model := new(M)
	err := b.executor.Run(ctx, func(idb bun.IDB) error {
		return idb.NewSelect().Model(model).Where("id = ?", id).Scan(ctx)
	})
	if err != nil {
		return nil, transformNoRows(err)
	}
	return model, nil
}

func (b *baseRepo[M]) UpdateByID(ctx context.Context, id uuid.UUID, model *M) error {
	return b.executor.Run(ctx, func(idb bun.IDB) error {
		result, err := idb.NewUpdate().Model(model).Where("id = ?", id).Exec(ctx)
		if err != nil {
			return err
		}
		return transformRowsAffected(result)
	})
}

func (b *baseRepo[M]) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return b.executor.Run(ctx, func(idb bun.IDB) error {
		result, err := idb.NewDelete().Model(new(M)).Where("id = ?", id).Exec(ctx)
		if err != nil {
			return err
		}
		return transformRowsAffected(result)
	})
}

func transformNoRows(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}

func transformRowsAffected(result sql.Result) error {
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
