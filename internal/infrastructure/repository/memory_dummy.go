// Package repository implements driven adapters for persistence.
package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/fermin/gophercraft/internal/domain/model"
	"github.com/fermin/gophercraft/internal/domain/port"
)

var _ port.DummyRepository = (*MemoryDummyRepository)(nil)

// MemoryDummyRepository is an in-memory adapter for DummyRepository.
type MemoryDummyRepository struct {
	mu   sync.RWMutex
	data map[string]model.Dummy
}

// NewMemoryDummyRepository returns an empty in-memory store.
func NewMemoryDummyRepository() *MemoryDummyRepository {
	return &MemoryDummyRepository{
		data: make(map[string]model.Dummy),
	}
}

// Save persists or overwrites a Dummy by ID.
func (r *MemoryDummyRepository) Save(_ context.Context, e model.Dummy) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[e.ID().String()] = e
	return nil
}

// FindByID returns a Dummy by ID.
func (r *MemoryDummyRepository) FindByID(_ context.Context, id model.DummyID) (model.Dummy, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.data[id.String()]
	if !ok {
		return model.Dummy{}, errors.New("dummy not found")
	}
	return e, nil
}
