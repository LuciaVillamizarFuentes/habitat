package unit

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Service contiene los casos de uso de unidades.
type Service struct {
	repo Repository
	now  func() time.Time // inyectable para tests deterministas
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo, now: time.Now}
}

// Register crea un nuevo inquilino tras validar los datos.
func (s *Service) Register(ctx context.Context, in CreateUnit) (Unit, error) {
	if err := in.Validate(); err != nil {
		return Unit{}, err
	}

	if _, err := s.repo.GetByCode(ctx, in.Code); err == nil {
		return Unit{}, ErrDuplicateCode
	} else if !errors.Is(err, ErrNotFound) {
		return Unit{}, fmt.Errorf("check code: %w", err)
	}

	now := s.now().UTC()
	t := Unit{
		ID:          newID(),
		Code:        in.Code,
		Kind:        in.Kind,
		Floor:       in.Floor,
		AreaM2:      in.AreaM2,
		Coefficient: in.Coefficient,
		CreatedAt:   now,
	}
	if err := s.repo.Create(ctx, t); err != nil {
		return Unit{}, fmt.Errorf("create unit: %w", err)
	}

	return t, nil
}

func (s *Service) Get(ctx context.Context, id string) (Unit, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]Unit, error) {
	return s.repo.List(ctx)
}

// newID genera un identificador aleatorio.
// TODO(T-06): reemplazar por UUIDv7 (ordenable por tiempo) y justificar la decisión en un ADR.
func newID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
