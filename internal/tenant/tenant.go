// Package tenant gestiona a los inquilinos (residentes) del edificio.
//
// Es el módulo de REFERENCIA del template: está completo (dominio, servicio,
// repositorio en memoria, handlers HTTP y tests). Úsalo como guía para
// construir los demás módulos.
package tenant

import (
	"errors"
	"net/mail"
	"strings"
	"time"
)

// Errores de dominio. Los handlers los traducen a códigos HTTP.
var (
	ErrNotFound      = errors.New("tenant not found")
	ErrInvalidName   = errors.New("name is required")
	ErrInvalidEmail  = errors.New("email is invalid")
	ErrDuplicateMail = errors.New("email already registered")
)

// Tenant es un residente del edificio.
type Tenant struct {
	ID        string    `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone,omitempty"`
	UnitID    string    `json:"unit_id,omitempty"` // TODO(T-05): relacionar con unit.Unit
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateInput son los datos necesarios para registrar un inquilino.
type CreateInput struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	UnitID   string `json:"unit_id"`
}

// Validate aplica las reglas de negocio mínimas.
func (in *CreateInput) Validate() error {
	in.FullName = strings.TrimSpace(in.FullName)
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	if in.FullName == "" {
		return ErrInvalidName
	}
	if _, err := mail.ParseAddress(in.Email); err != nil {
		return ErrInvalidEmail
	}
	return nil
}
