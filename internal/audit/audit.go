// Package audit guarda la HISTORIA del edificio: quién hizo qué y cuándo.
// Es un log append-only de eventos de dominio.
//
// ESTADO: solo dominio. Implementación → T-15, T-19, T-22 del ROADMAP.
package audit

import (
	"encoding/json"
	"time"
)

// Event es un hecho inmutable que ocurrió en el sistema.
type Event struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`         // "tenant.registered", "reservation.confirmed"...
	AggregateID string          `json:"aggregate_id"` // ID de la entidad afectada
	ActorID     string          `json:"actor_id"`     // quién lo causó
	Payload     json.RawMessage `json:"payload"`
	OccurredAt  time.Time       `json:"occurred_at"`
}
