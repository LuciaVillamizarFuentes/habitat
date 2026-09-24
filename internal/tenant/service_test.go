package tenant

import (
	"context"
	"errors"
	"testing"
)

// Ejemplo de table-driven test: el patrón idiomático en Go.
func TestService_Register(t *testing.T) {
	tests := []struct {
		name    string
		seed    []CreateInput
		input   CreateInput
		wantErr error
	}{
		{
			name:  "ok",
			input: CreateInput{FullName: "Ana Gómez", Email: "ana@example.com"},
		},
		{
			name:    "nombre vacío",
			input:   CreateInput{FullName: "   ", Email: "ana@example.com"},
			wantErr: ErrInvalidName,
		},
		{
			name:    "email inválido",
			input:   CreateInput{FullName: "Ana", Email: "no-es-email"},
			wantErr: ErrInvalidEmail,
		},
		{
			name:    "email duplicado (case insensitive)",
			seed:    []CreateInput{{FullName: "Ana", Email: "ana@example.com"}},
			input:   CreateInput{FullName: "Otra Ana", Email: "ANA@example.com"},
			wantErr: ErrDuplicateMail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(NewMemoryRepository())
			ctx := context.Background()
			for _, s := range tt.seed {
				if _, err := svc.Register(ctx, s); err != nil {
					t.Fatalf("seed: %v", err)
				}
			}

			got, err := svc.Register(ctx, tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && (got.ID == "" || !got.Active) {
				t.Fatalf("tenant inválido: %+v", got)
			}
		})
	}
}
