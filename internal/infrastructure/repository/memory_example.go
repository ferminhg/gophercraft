package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/fermin/gophercraft/internal/domain/model"
	"github.com/fermin/gophercraft/internal/domain/port"
)

var _ port.ExampleRepository = (*MemoryExampleRepository)(nil)

// MemoryExampleRepository is an in-memory adapter for ExampleRepository.
type MemoryExampleRepository struct {
	mu   sync.RWMutex
	data map[string]model.Example
}

// NewMemoryExampleRepository returns an empty in-memory store.
func NewMemoryExampleRepository() *MemoryExampleRepository {
	return &MemoryExampleRepository{
		data: make(map[string]model.Example),
	}
}

// Save persists or overwrites an Example by ID.
func (r *MemoryExampleRepository) Save(_ context.Context, e model.Example) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[e.ID] = e
	return nil
}

// FindByID returns an Example by ID.
func (r *MemoryExampleRepository) FindByID(_ context.Context, id string) (model.Example, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.data[id]
	if !ok {
		return model.Example{}, errors.New("example not found")
	}
	return e, nil
}
