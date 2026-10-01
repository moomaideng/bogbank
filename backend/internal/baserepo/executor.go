package baserepo

import (
	"context"

	"github.com/uptrace/bun"
)

type dbContextKey struct{}

// Executor runs fn against the transaction on ctx, or against the raw DB when
// there is none. Repositories take an Executor so callers do not thread a
// transaction through every method.
type Executor interface {
	Run(ctx context.Context, fn func(idb bun.IDB) error) error
}

// Transactioner stores a bun transaction on ctx for nested Executor.Run calls.
type Transactioner interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type executorImpl struct{ db *bun.DB }

type transactionerImpl struct{ db *bun.DB }

// NewExecutor uses db when ctx carries no transaction.
func NewExecutor(db *bun.DB) Executor {
	return &executorImpl{db: db}
}

func (e *executorImpl) Run(ctx context.Context, fn func(idb bun.IDB) error) error {
	idb, ok := ctx.Value(dbContextKey{}).(bun.IDB)
	if !ok {
		idb = e.db
	}
	return fn(idb)
}

// NewTransactioner starts transactions on db.
func NewTransactioner(db *bun.DB) Transactioner {
	return &transactionerImpl{db: db}
}

func (t *transactionerImpl) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		return fn(context.WithValue(ctx, dbContextKey{}, tx))
	})
}
