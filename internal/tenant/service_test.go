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

func strPtr(s string) *string { return &s }
func boolPtr(b bool) *bool    { return &b }

func TestService_Update(t *testing.T) {
	tests := []struct {
		name      string
		otherSeed []CreateInput // otros inquilinos ya registrados, para probar colisión de email
		input     UpdateInput
		wantErr   error
	}{
		{
			name:  "ok nombre",
			input: UpdateInput{FullName: strPtr("Ana Actualizada")},
		},
		{
			name:  "ok email",
			input: UpdateInput{Email: strPtr("nuevo@example.com")},
		},
		{
			name:  "ok inactivar",
			input: UpdateInput{Active: boolPtr(false)},
		},
		{
			name:    "nombre vacío",
			input:   UpdateInput{FullName: strPtr("   ")},
			wantErr: ErrInvalidName,
		},
		{
			name:    "email inválido",
			input:   UpdateInput{Email: strPtr("no-es-email")},
			wantErr: ErrInvalidEmail,
		},
		{
			name:  "email igual al actual no se considera duplicado",
			input: UpdateInput{Email: strPtr("ana@example.com")},
		},
		{
			name:  "email igual al actual con distinto case no se considera duplicado",
			input: UpdateInput{Email: strPtr("ANA@example.com")},
		},
		{
			name:      "email de otro inquilino es duplicado",
			otherSeed: []CreateInput{{FullName: "Otra", Email: "otra@example.com"}},
			input:     UpdateInput{Email: strPtr("otra@example.com")},
			wantErr:   ErrDuplicateMail,
		},
		{
			name:      "email de otro inquilino es duplicado (case insensitive)",
			otherSeed: []CreateInput{{FullName: "Otra", Email: "otra@example.com"}},
			input:     UpdateInput{Email: strPtr("OTRA@example.com")},
			wantErr:   ErrDuplicateMail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(NewMemoryRepository())
			ctx := context.Background()
			for _, s := range tt.otherSeed {
				if _, err := svc.Register(ctx, s); err != nil {
					t.Fatalf("seed otro inquilino: %v", err)
				}
			}
			created, err := svc.Register(ctx, CreateInput{FullName: "Ana Gómez", Email: "ana@example.com"})
			if err != nil {
				t.Fatalf("seed: %v", err)
			}

			got, err := svc.Update(ctx, created.ID, tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}

			if !got.UpdatedAt.After(created.UpdatedAt) && !got.UpdatedAt.Equal(created.UpdatedAt) {
				t.Fatalf("UpdatedAt = %v, want >= %v", got.UpdatedAt, created.UpdatedAt)
			}
			if tt.input.FullName != nil && got.FullName != *tt.input.FullName {
				t.Fatalf("FullName = %q, want %q", got.FullName, *tt.input.FullName)
			}
			if tt.input.Email != nil && got.Email != *tt.input.Email {
				t.Fatalf("Email = %q, want %q", got.Email, *tt.input.Email)
			}
			if tt.input.Active != nil && got.Active != *tt.input.Active {
				t.Fatalf("Active = %v, want %v", got.Active, *tt.input.Active)
			}
		})
	}

	t.Run("no existe", func(t *testing.T) {
		svc := NewService(NewMemoryRepository())
		_, err := svc.Update(context.Background(), "no-existe", UpdateInput{FullName: strPtr("X")})
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("err = %v, want %v", err, ErrNotFound)
		}
	})
}

func TestService_Delete(t *testing.T) {
	svc := NewService(NewMemoryRepository())
	ctx := context.Background()
	created, err := svc.Register(ctx, CreateInput{FullName: "Ana Gómez", Email: "ana@example.com"})
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	if err := svc.Delete(ctx, created.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	got, err := svc.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Active {
		t.Fatalf("Active = true, want false tras soft delete")
	}
	if !got.UpdatedAt.After(created.UpdatedAt) && !got.UpdatedAt.Equal(created.UpdatedAt) {
		t.Fatalf("UpdatedAt = %v, want >= %v", got.UpdatedAt, created.UpdatedAt)
	}

	t.Run("no existe", func(t *testing.T) {
		if err := svc.Delete(ctx, "no-existe"); !errors.Is(err, ErrNotFound) {
			t.Fatalf("err = %v, want %v", err, ErrNotFound)
		}
	})
}
