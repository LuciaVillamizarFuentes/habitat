// Package unit modela las unidades del edificio: apartamentos, locales, parqueaderos, depósitos.
//
// ESTADO: solo dominio. Implementación → tareas T-01, T-02, T-05 del ROADMAP.
package unit

import (
	"errors"
	"strings"
	"time"
)

type Kind string

const (
	KindApartment Kind = "apartment"
	KindParking   Kind = "parking"
	KindStorage   Kind = "storage"
	KindCommerce  Kind = "commerce"
)

var validKinds = map[Kind]bool{
	KindApartment: true,
	KindParking:   true,
	KindStorage:   true,
	KindCommerce:  true,
}

var (
	ErrInvalidCode        = errors.New("code is required")
	ErrInvalidKind        = errors.New("kind is required")
	ErrInvalidFloor       = errors.New("floor is invalid")
	ErrInvalidArea        = errors.New("area is invalid")
	ErrInvalidCoefficient = errors.New("coefficient is invalid")
	ErrNotFound           = errors.New("unit not found")
	ErrDuplicateCode      = errors.New("code already exists")
)

// Unit es una unidad física del edificio.
type Unit struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"` // ej: "T1-502"
	Kind        Kind      `json:"kind"`
	Floor       int       `json:"floor"`
	AreaM2      float64   `json:"area_m2"`
	Coefficient float64   `json:"coefficient"` // coeficiente de copropiedad (Ley 675 de 2001)
	CreatedAt   time.Time `json:"created_at"`
}

type CreateUnit struct {
	Code        string  `json:"code"`
	Kind        Kind    `json:"kind"`
	Floor       int     `json:"floor"`
	AreaM2      float64 `json:"area_m2"`
	Coefficient float64 `json:"coefficient"` // coeficiente de copropiedad (Ley 675 de 2001)
}

// TODO(T-02): Repository + Service + Handler siguiendo el módulo tenant..
func (in *CreateUnit) Validate() error {
	in.Code = strings.TrimSpace(in.Code)
	if in.Code == "" {
		return ErrInvalidCode
	}
	in.Kind = Kind(strings.TrimSpace(string(in.Kind)))
	if !validKinds[in.Kind] {
		return ErrInvalidKind
	}
	if in.Floor < -5 {
		return ErrInvalidFloor
	}
	if in.AreaM2 <= 0 {
		return ErrInvalidArea
	}
	if in.Coefficient <= 0 || in.Coefficient > 1 {
		return ErrInvalidCoefficient
	}
	return nil
}
