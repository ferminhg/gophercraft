// Package command contains write-side application use cases.
package command

import (
	"context"

	"github.com/fermin/gophercraft/internal/domain/model"
	"github.com/fermin/gophercraft/internal/domain/port"
)

// CreateDummyHandler handles create use cases (command side).
type CreateDummyHandler struct {
	repo port.DummyRepository
}

// NewCreateDummyHandler constructs a handler wired to the repository port.
func NewCreateDummyHandler(repo port.DummyRepository) *CreateDummyHandler {
	return &CreateDummyHandler{repo: repo}
}

// Handle persists a new Dummy produced by the domain constructors.
func (h *CreateDummyHandler) Handle(ctx context.Context, e model.Dummy) error {
	return h.repo.Save(ctx, e)
}
