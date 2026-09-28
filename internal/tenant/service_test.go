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

// seedForList registra 3 inquilinos (Ana, Carlos, Beatriz, en ese orden de
// creación) y desactiva a Beatriz, para poder probar cada filtro por separado.
func seedForList(t *testing.T, svc *Service) []Tenant {
	t.Helper()
	ctx := context.Background()
	inputs := []CreateInput{
		{FullName: "Ana Gómez", Email: "ana@example.com", Phone: "111", UnitID: "U1"},
		{FullName: "Carlos Ruiz", Email: "carlos@example.com", Phone: "222", UnitID: "U2"},
		{FullName: "Beatriz Soto", Email: "beatriz@example.com", Phone: "333", UnitID: "U1"},
	}
	out := make([]Tenant, 0, len(inputs))
	for _, in := range inputs {
		created, err := svc.Register(ctx, in)
		if err != nil {
			t.Fatalf("seed: %v", err)
		}
		out = append(out, created)
	}
	if err := svc.Delete(ctx, out[2].ID); err != nil {
		t.Fatalf("desactivar Beatriz: %v", err)
	}
	out[2].Active = false
	return out
}

func TestService_List_Filters(t *testing.T) {
	tests := []struct {
		name       string
		filter     FilterTenants
		wantEmails []string // en el orden esperado (por CreatedAt asc)
	}{
		{
			name:       "sin filtro devuelve todos",
			filter:     FilterTenants{},
			wantEmails: []string{"ana@example.com", "carlos@example.com", "beatriz@example.com"},
		},
		{
			name:       "full_name exacto",
			filter:     FilterTenants{FullName: strPtr("Ana Gómez")},
			wantEmails: []string{"ana@example.com"},
		},
		{
			name:       "full_name no coincide",
			filter:     FilterTenants{FullName: strPtr("No Existe")},
			wantEmails: nil,
		},
		{
			name:       "email case insensitive",
			filter:     FilterTenants{Email: strPtr("ANA@Example.com")},
			wantEmails: []string{"ana@example.com"},
		},
		{
			name:       "phone exacto",
			filter:     FilterTenants{Phone: strPtr("222")},
			wantEmails: []string{"carlos@example.com"},
		},
		{
			name:       "unit_id con varios resultados",
			filter:     FilterTenants{UnitID: strPtr("U1")},
			wantEmails: []string{"ana@example.com", "beatriz@example.com"},
		},
		{
			name:       "active true excluye inactivos",
			filter:     FilterTenants{Active: boolPtr(true)},
			wantEmails: []string{"ana@example.com", "carlos@example.com"},
		},
		{
			name:       "active false solo inactivos",
			filter:     FilterTenants{Active: boolPtr(false)},
			wantEmails: []string{"beatriz@example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(NewMemoryRepository())
			seedForList(t, svc)

			got, err := svc.List(context.Background(), 10, 1, tt.filter)
			if err != nil {
				t.Fatalf("List: %v", err)
			}

			gotEmails := make([]string, len(got))
			for i, tn := range got {
				gotEmails[i] = tn.Email
			}
			if len(gotEmails) != len(tt.wantEmails) {
				t.Fatalf("emails = %v, want %v", gotEmails, tt.wantEmails)
			}
			for i := range gotEmails {
				if gotEmails[i] != tt.wantEmails[i] {
					t.Fatalf("emails = %v, want %v", gotEmails, tt.wantEmails)
				}
			}
		})
	}
}

func TestService_List_Pagination(t *testing.T) {
	svc := NewService(NewMemoryRepository())
	seeded := seedForList(t, svc) // 3 inquilinos: Ana, Carlos, Beatriz (en ese orden)
	ctx := context.Background()

	tests := []struct {
		name       string
		limit      uint64
		page       uint64
		wantEmails []string
	}{
		{name: "página 1", limit: 2, page: 1, wantEmails: []string{seeded[0].Email, seeded[1].Email}},
		{name: "página 2", limit: 2, page: 2, wantEmails: []string{seeded[2].Email}},
		{name: "página fuera de rango", limit: 2, page: 3, wantEmails: []string{}},
		{name: "page en 0 se trata como la primera página", limit: 2, page: 0, wantEmails: []string{seeded[0].Email, seeded[1].Email}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := svc.List(ctx, tt.limit, tt.page, FilterTenants{})
			if err != nil {
				t.Fatalf("List: %v", err)
			}
			gotEmails := make([]string, len(got))
			for i, tn := range got {
				gotEmails[i] = tn.Email
			}
			if len(gotEmails) != len(tt.wantEmails) {
				t.Fatalf("emails = %v, want %v", gotEmails, tt.wantEmails)
			}
			for i := range gotEmails {
				if gotEmails[i] != tt.wantEmails[i] {
					t.Fatalf("emails = %v, want %v", gotEmails, tt.wantEmails)
				}
			}
		})
	}
}
