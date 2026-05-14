package command

import (
	"context"

	"github.com/fermin/gophercraft/internal/domain/model"
	"github.com/fermin/gophercraft/internal/domain/port"
)

// CreateExampleHandler handles create use cases (command side).
type CreateExampleHandler struct {
	repo port.ExampleRepository
}

// NewCreateExampleHandler constructs a handler wired to the repository port.
func NewCreateExampleHandler(repo port.ExampleRepository) *CreateExampleHandler {
	return &CreateExampleHandler{repo: repo}
}

// Handle persists a new Example when the model is valid.
func (h *CreateExampleHandler) Handle(ctx context.Context, e model.Example) error {
	if !e.IsValid() {
		return ErrInvalidExample
	}
	return h.repo.Save(ctx, e)
}
