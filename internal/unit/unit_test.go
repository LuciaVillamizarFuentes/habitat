package unit

import (
	"errors"
	"testing"
)

func TestCreateUnit_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   CreateUnit
		wantErr error
	}{
		{
			name: "ok apartment",
			input: CreateUnit{
				Code:        "T1-502",
				Kind:        KindApartment,
				Floor:       5,
				AreaM2:      80,
				Coefficient: 0.015,
			},
		},
		{
			name: "ok parking en sótano",
			input: CreateUnit{
				Code:        "P-12",
				Kind:        KindParking,
				Floor:       -2,
				AreaM2:      12,
				Coefficient: 0.001,
			},
		},
		{
			name: "ok coeficiente máximo permitido",
			input: CreateUnit{
				Code:        "L-1",
				Kind:        KindCommerce,
				Floor:       1,
				AreaM2:      40,
				Coefficient: 1,
			},
		},
		{
			name: "código vacío",
			input: CreateUnit{
				Code:        "   ",
				Kind:        KindApartment,
				Floor:       1,
				AreaM2:      80,
				Coefficient: 0.01,
			},
			wantErr: ErrInvalidCode,
		},
		{
			name: "kind vacío",
			input: CreateUnit{
				Code:        "T1-502",
				Kind:        "",
				Floor:       1,
				AreaM2:      80,
				Coefficient: 0.01,
			},
			wantErr: ErrInvalidKind,
		},
		{
			name: "kind desconocido",
			input: CreateUnit{
				Code:        "T1-502",
				Kind:        "garage",
				Floor:       1,
				AreaM2:      80,
				Coefficient: 0.01,
			},
			wantErr: ErrInvalidKind,
		},
		{
			name: "floor por debajo del mínimo",
			input: CreateUnit{
				Code:        "S-1",
				Kind:        KindStorage,
				Floor:       -6,
				AreaM2:      5,
				Coefficient: 0.001,
			},
			wantErr: ErrInvalidFloor,
		},
		{
			name: "area cero",
			input: CreateUnit{
				Code:        "T1-502",
				Kind:        KindApartment,
				Floor:       1,
				AreaM2:      0,
				Coefficient: 0.01,
			},
			wantErr: ErrInvalidArea,
		},
		{
			name: "area negativa",
			input: CreateUnit{
				Code:        "T1-502",
				Kind:        KindApartment,
				Floor:       1,
				AreaM2:      -10,
				Coefficient: 0.01,
			},
			wantErr: ErrInvalidArea,
		},
		{
			name: "coeficiente cero",
			input: CreateUnit{
				Code:        "T1-502",
				Kind:        KindApartment,
				Floor:       1,
				AreaM2:      80,
				Coefficient: 0,
			},
			wantErr: ErrInvalidCoefficient,
		},
		{
			name: "coeficiente negativo",
			input: CreateUnit{
				Code:        "T1-502",
				Kind:        KindApartment,
				Floor:       1,
				AreaM2:      80,
				Coefficient: -0.01,
			},
			wantErr: ErrInvalidCoefficient,
		},
		{
			name: "coeficiente mayor a 1",
			input: CreateUnit{
				Code:        "T1-502",
				Kind:        KindApartment,
				Floor:       1,
				AreaM2:      80,
				Coefficient: 1.01,
			},
			wantErr: ErrInvalidCoefficient,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateUnit_Validate_TrimsCodeAndKind(t *testing.T) {
	in := CreateUnit{
		Code:        "  T1-502  ",
		Kind:        "  apartment  ",
		Floor:       1,
		AreaM2:      80,
		Coefficient: 0.01,
	}

	if err := in.Validate(); err != nil {
		t.Fatalf("err = %v, want nil", err)
	}
	if in.Code != "T1-502" {
		t.Fatalf("Code = %q, want %q", in.Code, "T1-502")
	}
	if in.Kind != KindApartment {
		t.Fatalf("Kind = %q, want %q", in.Kind, KindApartment)
	}
}
