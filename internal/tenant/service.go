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

func (s *Service) List(ctx context.Context, limit uint64, page uint64, filter FilterTenants) ([]Tenant, error) {
	return s.repo.List(ctx, limit, page, filter)
}

func (s *Service) Update(ctx context.Context, id string, in UpdateInput) (Tenant, error) {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return Tenant{}, fmt.Errorf("get tenant: %w", err)
	}

	if err := in.Validate(); err != nil {
		return Tenant{}, err
	}

	if in.FullName != nil {
		t.FullName = *in.FullName
	}
	if in.Email != nil && t.Email != *in.Email {
		t.Email = *in.Email
		if _, err := s.repo.GetByEmail(ctx, t.Email); err == nil {
			return Tenant{}, ErrDuplicateMail
		} else if !errors.Is(err, ErrNotFound) {
			return Tenant{}, fmt.Errorf("check email: %w", err)
		}
	}
	if in.Phone != nil {
		t.Phone = *in.Phone
	}
	if in.UnitID != nil {
		t.UnitID = *in.UnitID
	}
	if in.Active != nil {
		t.Active = *in.Active
	}
	t.UpdatedAt = s.now().UTC()

	if err := s.repo.Update(ctx, t); err != nil {
		return Tenant{}, fmt.Errorf("update tenant: %w", err)
	}
	return t, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	t, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get tenant: %w", err)
	}
	t.Active = false
	t.UpdatedAt = s.now().UTC()
	if err := s.repo.Update(ctx, t); err != nil {
		return fmt.Errorf("soft delete tenant: %w", err)
	}
	return nil
}

// newID genera un identificador aleatorio.
// TODO(T-06): reemplazar por UUIDv7 (ordenable por tiempo) y justificar la decisión en un ADR.
func newID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
