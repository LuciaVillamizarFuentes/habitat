package tenant

import (
	"context"
	"sort"
	"sync"
)

// Repository define la persistencia de inquilinos.
// El servicio depende de esta interfaz, no de una implementación concreta.
type Repository interface {
	Create(ctx context.Context, t Tenant) error
	GetByID(ctx context.Context, id string) (Tenant, error)
	GetByEmail(ctx context.Context, email string) (Tenant, error)
	List(ctx context.Context) ([]Tenant, error)
}

// MemoryRepository es una implementación en memoria, útil para desarrollo y tests.
// TODO(T-08): crear PostgresRepository que implemente la misma interfaz.
type MemoryRepository struct {
	mu   sync.RWMutex
	data map[string]Tenant
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{data: make(map[string]Tenant)}
}

func (r *MemoryRepository) Create(_ context.Context, t Tenant) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.data {
		if existing.Email == t.Email {
			return ErrDuplicateMail
		}
	}
	r.data[t.ID] = t
	return nil
}

func (r *MemoryRepository) GetByID(_ context.Context, id string) (Tenant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.data[id]
	if !ok {
		return Tenant{}, ErrNotFound
	}
	return t, nil
}

func (r *MemoryRepository) GetByEmail(_ context.Context, email string) (Tenant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.data {
		if t.Email == email {
			return t, nil
		}
	}
	return Tenant{}, ErrNotFound
}

// List devuelve todos los inquilinos ordenados por fecha de creación.
// TODO(T-04): añadir paginación y filtros (active, unit_id).
func (r *MemoryRepository) List(_ context.Context) ([]Tenant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Tenant, 0, len(r.data))
	for _, t := range r.data {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out, nil
}
