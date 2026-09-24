// Package unit modela las unidades del edificio: apartamentos, locales, parqueaderos, depósitos.
//
// ESTADO: solo dominio. Implementación → tareas T-01, T-02, T-05 del ROADMAP.
package unit

import "time"

type Kind string

const (
	KindApartment Kind = "apartment"
	KindParking   Kind = "parking"
	KindStorage   Kind = "storage"
	KindCommerce  Kind = "commerce"
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

// TODO(T-01): validar Kind, Code único, Floor >= -5, AreaM2 > 0, Coefficient en (0, 1].
// TODO(T-02): Repository + Service + Handler siguiendo el módulo tenant.
