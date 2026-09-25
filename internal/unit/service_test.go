package unit

import (
	"context"
	"errors"
	"testing"
)

func TestService_Register(t *testing.T) {
	tests := []struct {
		name    string
		seed    []CreateUnit
		input   CreateUnit
		wantErr error
	}{
		{
			name:  "ok",
			input: CreateUnit{Code: "T1-502", Kind: KindApartment, Floor: 5, AreaM2: 80, Coefficient: 0.015},
		},
		{
			name:    "código vacío",
			input:   CreateUnit{Code: "   ", Kind: KindApartment, Floor: 5, AreaM2: 80, Coefficient: 0.015},
			wantErr: ErrInvalidCode,
		},
		{
			name:    "kind inválido",
			input:   CreateUnit{Code: "T1-502", Kind: "garage", Floor: 5, AreaM2: 80, Coefficient: 0.015},
			wantErr: ErrInvalidKind,
		},
		{
			name:    "floor inválido",
			input:   CreateUnit{Code: "T1-502", Kind: KindApartment, Floor: -6, AreaM2: 80, Coefficient: 0.015},
			wantErr: ErrInvalidFloor,
		},
		{
			name:    "area inválida",
			input:   CreateUnit{Code: "T1-502", Kind: KindApartment, Floor: 5, AreaM2: 0, Coefficient: 0.015},
			wantErr: ErrInvalidArea,
		},
		{
			name:    "coeficiente inválido",
			input:   CreateUnit{Code: "T1-502", Kind: KindApartment, Floor: 5, AreaM2: 80, Coefficient: 0},
			wantErr: ErrInvalidCoefficient,
		},
		{
			name:    "código duplicado",
			seed:    []CreateUnit{{Code: "T1-502", Kind: KindApartment, Floor: 5, AreaM2: 80, Coefficient: 0.015}},
			input:   CreateUnit{Code: "T1-502", Kind: KindParking, Floor: -1, AreaM2: 12, Coefficient: 0.001},
			wantErr: ErrDuplicateCode,
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
			if tt.wantErr == nil && (got.ID == "" || got.CreatedAt.IsZero()) {
				t.Fatalf("unit inválida: %+v", got)
			}
		})
	}
}

func TestService_GetAndList(t *testing.T) {
	svc := NewService(NewMemoryRepository())
	ctx := context.Background()

	if _, err := svc.Get(ctx, "no-existe"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get err = %v, want %v", err, ErrNotFound)
	}

	created, err := svc.Register(ctx, CreateUnit{Code: "T1-502", Kind: KindApartment, Floor: 5, AreaM2: 80, Coefficient: 0.015})
	if err != nil {
		t.Fatalf("Register: %v", err)
	}

	got, err := svc.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Code != created.Code {
		t.Fatalf("Code = %q, want %q", got.Code, created.Code)
	}

	list, err := svc.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 1 || list[0].ID != created.ID {
		t.Fatalf("List = %+v, want [%+v]", list, created)
	}
}
