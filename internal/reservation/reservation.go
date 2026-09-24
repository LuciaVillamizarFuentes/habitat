// Package reservation gestiona las reservas de zonas comunes.
//
// Este es el CORAZÓN del reto técnico: evitar doble reserva bajo concurrencia.
// ESTADO: solo dominio. Implementación → T-10, T-17, T-21 del ROADMAP.
package reservation

import (
	"errors"
	"time"
)

var (
	ErrOverlap      = errors.New("the resource is already booked in that time range")
	ErrOutsideHours = errors.New("reservation outside resource opening hours")
	ErrInvalidRange = errors.New("end must be after start")
)

type Status string

const (
	StatusPending   Status = "pending" // esperando pago/aprobación
	StatusConfirmed Status = "confirmed"
	StatusCancelled Status = "cancelled"
)

type Reservation struct {
	ID         string    `json:"id"`
	ResourceID string    `json:"resource_id"`
	TenantID   string    `json:"tenant_id"`
	StartsAt   time.Time `json:"starts_at"`
	EndsAt     time.Time `json:"ends_at"`
	Status     Status    `json:"status"`
	Version    int       `json:"version"` // para optimistic locking (T-17)
	CreatedAt  time.Time `json:"created_at"`
}

// Overlaps indica si dos rangos [a, b) se solapan.
// TODO(T-10): escribir tests de borde (rangos contiguos NO se solapan).
func Overlaps(aStart, aEnd, bStart, bEnd time.Time) bool {
	return aStart.Before(bEnd) && bStart.Before(aEnd)
}
