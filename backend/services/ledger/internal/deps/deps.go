// Package deps builds the usecases the ledger transports call.
package deps

import (
	"github.com/uptrace/bun"

	"github.com/moomaideng/bogbank/backend/services/ledger/internal/adapter/postgres"
	"github.com/moomaideng/bogbank/backend/services/ledger/internal/usecase"
)

type Deps struct {
	Categories *usecase.CategoryUsecase
}

func New(db *bun.DB) Deps {
	return Deps{
		Categories: usecase.NewCategoryUsecase(postgres.NewCategoryRepo(db)),
	}
}
