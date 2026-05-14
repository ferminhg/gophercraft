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

// Handle persists a new Dummy when the model is valid.
func (h *CreateDummyHandler) Handle(ctx context.Context, e model.Dummy) error {
	if !e.IsValid() {
		return ErrInvalidDummy
	}
	return h.repo.Save(ctx, e)
}
