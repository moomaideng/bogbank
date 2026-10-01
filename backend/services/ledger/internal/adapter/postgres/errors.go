package postgres

import (
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/moomaideng/bogbank/backend/internal/baserepo"
	"github.com/moomaideng/bogbank/backend/services/ledger/internal/usecase"
)

func mapReadErr(err error) error {
	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, baserepo.ErrNotFound) {
		return usecase.ErrNotFound
	}
	return err
}

func mapWriteErr(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return usecase.ErrConflict
	}
	if errors.Is(err, baserepo.ErrNotFound) {
		return usecase.ErrNotFound
	}
	return err
}
