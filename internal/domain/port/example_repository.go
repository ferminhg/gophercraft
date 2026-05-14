package port

import (
	"context"

	"github.com/fermin/gophercraft/internal/domain/model"
)

// ExampleRepository is a driven port for Example persistence.
type ExampleRepository interface {
	Save(ctx context.Context, e model.Example) error
	FindByID(ctx context.Context, id string) (model.Example, error)
}
