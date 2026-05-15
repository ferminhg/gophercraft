// Package command contains write-side application use cases.
package command

import (
	"context"

	"github.com/fermin/gophercraft/internal/domain/model"
	"github.com/fermin/gophercraft/internal/domain/port"
)

// CreateDummyRequest carries primitives for creating a Dummy aggregate.
type CreateDummyRequest struct {
	Name string
	Type string
}

// CreateDummyHandler handles create Dummy use cases (command side).
type CreateDummyHandler struct {
	repo      port.DummyRepository
	uuidGen   port.UUIDGenerator
	clock     port.Clock
	publisher port.EventPublisher
}

// NewCreateDummyHandler constructs a handler wired to its ports.
func NewCreateDummyHandler(
	repo port.DummyRepository,
	uuidGen port.UUIDGenerator,
	clock port.Clock,
	publisher port.EventPublisher,
) *CreateDummyHandler {
	return &CreateDummyHandler{
		repo:      repo,
		uuidGen:   uuidGen,
		clock:     clock,
		publisher: publisher,
	}
}

// Handle persists a new Dummy produced by domain factory and publishes its domain events.
func (h *CreateDummyHandler) Handle(ctx context.Context, req CreateDummyRequest) error {
	id, err := model.NewDummyID(h.uuidGen.Generate())
	if err != nil {
		return err
	}
	name, err := model.NewDummyName(req.Name)
	if err != nil {
		return err
	}
	t, err := model.NewDummyType(req.Type)
	if err != nil {
		return err
	}
	createdAt, err := model.NewDummyCreatedAt(h.clock.Now())
	if err != nil {
		return err
	}
	d := model.CreateDummy(*id, *name, *t, *createdAt)
	if err := h.repo.Save(ctx, d); err != nil {
		return err
	}
	return h.publisher.Publish(ctx, d.PullDomainEvents()...)
}
