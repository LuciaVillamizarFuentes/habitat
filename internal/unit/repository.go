package unit

import (
	"context"
	"sort"
	"sync"
)

// Repository define la persistencia de unidades.
// El servicio depende de esta interfaz, no de una implementación concreta.
type Repository interface {
	Create(ctx context.Context, t Unit) error
	GetByID(ctx context.Context, id string) (Unit, error)
	GetByCode(ctx context.Context, code string) (Unit, error)
	List(ctx context.Context) ([]Unit, error)
}

// MemoryRepository es una implementación en memoria, útil para desarrollo y tests.

type MemoryRepository struct {
	mu   sync.RWMutex
	data map[string]Unit
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{data: make(map[string]Unit)}
}

func (r *MemoryRepository) Create(_ context.Context, t Unit) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.data {
		if existing.Code == t.Code {
			return ErrDuplicateCode
		}
	}
	r.data[t.ID] = t
	return nil
}

func (r *MemoryRepository) GetByID(_ context.Context, id string) (Unit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.data[id]
	if !ok {
		return Unit{}, ErrNotFound
	}
	return t, nil
}

func (r *MemoryRepository) GetByCode(_ context.Context, code string) (Unit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.data {
		if t.Code == code {
			return t, nil
		}
	}
	return Unit{}, ErrNotFound
}

// List devuelve todos los inquilinos ordenados por fecha de creación.
// TODO(T-04): añadir paginación y filtros (active, unit_id).
func (r *MemoryRepository) List(_ context.Context) ([]Unit, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Unit, 0, len(r.data))
	for _, t := range r.data {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out, nil
}
