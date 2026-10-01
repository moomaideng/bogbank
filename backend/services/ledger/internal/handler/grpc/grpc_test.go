package grpc

import (
	"errors"
	"fmt"
	"testing"

	"github.com/moomaideng/bogbank/backend/services/ledger/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestStatusErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		code codes.Code
		msg  string
	}{
		{
			name: "not found",
			err:  fmt.Errorf("wrapped: %w", usecase.ErrNotFound),
			code: codes.NotFound,
			msg:  "category not found",
		},
		{
			name: "conflict",
			err:  usecase.ErrConflict,
			code: codes.AlreadyExists,
			msg:  "category name already exists",
		},
		{
			name: "invalid argument",
			err:  usecase.ErrInvalidArgument,
			code: codes.InvalidArgument,
			msg:  "invalid argument",
		},
		{
			name: "unknown",
			err:  errors.New("driver: boom"),
			code: codes.Internal,
			msg:  "internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := statusErr(tt.err)
			st, ok := status.FromError(got)
			if !ok {
				t.Fatalf("got %v, want status error", got)
			}
			if st.Code() != tt.code {
				t.Fatalf("code = %v, want %v", st.Code(), tt.code)
			}
			if st.Message() != tt.msg {
				t.Fatalf("msg = %q, want %q", st.Message(), tt.msg)
			}
		})
	}
}
