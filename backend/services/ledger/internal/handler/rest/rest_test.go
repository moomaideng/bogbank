package rest

import (
	"errors"
	"fmt"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/moomaideng/bogbank/backend/services/ledger/internal/usecase"
)

func TestWriteErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		err    error
		status int
		msg    string
	}{
		{
			name:   "not found",
			err:    fmt.Errorf("wrapped: %w", usecase.ErrNotFound),
			status: 404,
			msg:    "category not found",
		},
		{
			name:   "conflict",
			err:    usecase.ErrConflict,
			status: 409,
			msg:    "category name already exists",
		},
		{
			name:   "invalid argument",
			err:    usecase.ErrInvalidArgument,
			status: 400,
			msg:    "invalid argument",
		},
		{
			name:   "unknown",
			err:    errors.New("driver: boom"),
			status: 500,
			msg:    "internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := writeErr(tt.err)
			se, ok := got.(huma.StatusError)
			if !ok {
				t.Fatalf("got %T, want huma.StatusError", got)
			}
			if se.GetStatus() != tt.status {
				t.Fatalf("status = %d, want %d", se.GetStatus(), tt.status)
			}
			if se.Error() != tt.msg {
				t.Fatalf("msg = %q, want %q", se.Error(), tt.msg)
			}
			if tt.status == 500 && se.Error() == "driver: boom" {
				t.Fatal("leaked internal error detail")
			}
		})
	}
}
