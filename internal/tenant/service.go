package tenant

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Service contiene los casos de uso de inquilinos.
type Service struct {
	repo Repository
	now  func() time.Time // inyectable para tests deterministas
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo, now: time.Now}
}

// Register crea un nuevo inquilino tras validar los datos.
func (s *Service) Register(ctx context.Context, in CreateInput) (Tenant, error) {
	if err := in.Validate(); err != nil {
		return Tenant{}, err
	}

	if _, err := s.repo.GetByEmail(ctx, in.Email); err == nil {
		return Tenant{}, ErrDuplicateMail
	} else if !errors.Is(err, ErrNotFound) {
		return Tenant{}, fmt.Errorf("check email: %w", err)
	}

	now := s.now().UTC()
	t := Tenant{
		ID:        newID(),
		FullName:  in.FullName,
		Email:     in.Email,
		Phone:     in.Phone,
		UnitID:    in.UnitID,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.Create(ctx, t); err != nil {
		return Tenant{}, fmt.Errorf("create tenant: %w", err)
	}
	// TODO(T-15): registrar el evento "tenant.registered" en el historial (audit).
	return t, nil
}

func (s *Service) Get(ctx context.Context, id string) (Tenant, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]Tenant, error) {
	return s.repo.List(ctx)
}

// newID genera un identificador aleatorio.
// TODO(T-06): reemplazar por UUIDv7 (ordenable por tiempo) y justificar la decisión en un ADR.
func newID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
