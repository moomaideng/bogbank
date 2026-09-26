package postgres

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/moomaideng/bogbank/internal/baserepo"
	"github.com/moomaideng/bogbank/services/ledger/internal/usecase"
)

func TestMapWriteErrConflict(t *testing.T) {
	t.Parallel()

	err := mapWriteErr(&pgconn.PgError{Code: "23505"})
	if !errors.Is(err, usecase.ErrConflict) {
		t.Fatalf("got %v, want ErrConflict", err)
	}
}

func TestMapWriteErrNotFound(t *testing.T) {
	t.Parallel()

	err := mapWriteErr(baserepo.ErrNotFound)
	if !errors.Is(err, usecase.ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestMapReadErr(t *testing.T) {
	t.Parallel()

	tests := []error{sql.ErrNoRows, baserepo.ErrNotFound}
	for _, in := range tests {
		err := mapReadErr(in)
		if !errors.Is(err, usecase.ErrNotFound) {
			t.Fatalf("mapReadErr(%v) = %v, want ErrNotFound", in, err)
		}
	}
}
