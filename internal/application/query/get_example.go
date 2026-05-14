package query

import (
	"context"

	"github.com/fermin/gophercraft/internal/domain/model"
	"github.com/fermin/gophercraft/internal/domain/port"
)

// GetExampleHandler handles read use cases (query side).
type GetExampleHandler struct {
	repo port.ExampleRepository
}

// NewGetExampleHandler constructs a handler wired to the repository port.
func NewGetExampleHandler(repo port.ExampleRepository) *GetExampleHandler {
	return &GetExampleHandler{repo: repo}
}

// Handle returns an Example by identifier.
func (h *GetExampleHandler) Handle(ctx context.Context, id string) (model.Example, error) {
	return h.repo.FindByID(ctx, id)
}
