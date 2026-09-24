// Package ticket gestiona las tareas/solicitudes del edificio: mantenimiento,
// quejas, PQRS, trabajos programados.
//
// ESTADO: solo dominio. Implementación → T-13, T-14 del ROADMAP.
package ticket

import "time"

type Status string

const (
	StatusOpen       Status = "open"
	StatusAssigned   Status = "assigned"
	StatusInProgress Status = "in_progress"
	StatusResolved   Status = "resolved"
	StatusClosed     Status = "closed"
)

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
	PriorityUrgent Priority = "urgent" // ej: fuga de agua, ascensor detenido
)

type Ticket struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Category    string     `json:"category"` // plomería, eléctrico, aseo, seguridad...
	Priority    Priority   `json:"priority"`
	Status      Status     `json:"status"`
	ReportedBy  string     `json:"reported_by"` // tenant ID
	AssignedTo  string     `json:"assigned_to,omitempty"`
	UnitID      string     `json:"unit_id,omitempty"`
	DueAt       *time.Time `json:"due_at,omitempty"` // SLA según prioridad
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TODO(T-13): implementar la máquina de estados. Transiciones válidas:
//   open -> assigned -> in_progress -> resolved -> closed
//   resolved -> in_progress (reabierto)
//   cualquier estado -> closed solo por un admin
