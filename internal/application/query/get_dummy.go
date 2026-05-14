package query

import (
	"context"

	"github.com/fermin/gophercraft/internal/domain/model"
	"github.com/fermin/gophercraft/internal/domain/port"
)

// GetDummyHandler handles read use cases (query side).
type GetDummyHandler struct {
	repo port.DummyRepository
}

// NewGetDummyHandler constructs a handler wired to the repository port.
func NewGetDummyHandler(repo port.DummyRepository) *GetDummyHandler {
	return &GetDummyHandler{repo: repo}
}

// Handle returns a Dummy by identifier.
func (h *GetDummyHandler) Handle(ctx context.Context, id string) (model.Dummy, error) {
	return h.repo.FindByID(ctx, id)
}
