// Package resource modela las zonas comunes y recursos reservables:
// salón social, BBQ, gimnasio, cancha, sala de juntas...
//
// ESTADO: solo dominio. Implementación → tarea T-09 del ROADMAP.
package resource

import "time"

// Resource es una zona común o recurso compartido.
type Resource struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	Capacity      int           `json:"capacity"`
	OpensAt       string        `json:"opens_at"`  // "08:00" (hora local del edificio)
	ClosesAt      string        `json:"closes_at"` // "22:00"
	MaxDuration   time.Duration `json:"max_duration"`
	RequiresFee   bool          `json:"requires_fee"`
	FeeCents      int64         `json:"fee_cents"` // dinero SIEMPRE en enteros, nunca float
	Active        bool          `json:"active"`
	MaintenanceAt *time.Time    `json:"maintenance_at,omitempty"`
}
