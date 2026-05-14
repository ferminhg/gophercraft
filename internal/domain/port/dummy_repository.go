package port

import (
	"context"

	"github.com/fermin/gophercraft/internal/domain/model"
)

// DummyRepository is a driven port for Dummy persistence.
type DummyRepository interface {
	Save(ctx context.Context, e model.Dummy) error
	FindByID(ctx context.Context, id string) (model.Dummy, error)
}
