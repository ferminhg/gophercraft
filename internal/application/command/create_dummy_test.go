package command_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/fermin/gophercraft/internal/application/command"
	"github.com/fermin/gophercraft/internal/domain/model"
	"github.com/fermin/gophercraft/internal/infrastructure/repository"
)

const dummyTestUUID = "550e8400-e29b-41d4-a716-446655440000"

type fixedUUIDGen struct{ value string }

func (g fixedUUIDGen) Generate() string { return g.value }

type saveFailsRepo struct{}

func (saveFailsRepo) Save(context.Context, *model.Dummy) error {
	return errors.New("save failed")
}

func (saveFailsRepo) FindByID(context.Context, model.DummyID) (model.Dummy, error) {
	return model.Dummy{}, errors.New("not found")
}

func TestCreateDummyHandler_Handle_PersistsAndPublishes(t *testing.T) {
	t.Parallel()

	repo := repository.NewMemoryDummyRepository()
	pub := &recordingPublisher{}
	now := time.Date(2026, 5, 14, 15, 30, 0, 0, time.UTC)

	h := command.NewCreateDummyHandler(
		repo,
		fixedUUIDGen{dummyTestUUID},
		fixedClock{instant: now},
		pub,
	)

	err := h.Handle(context.Background(), command.CreateDummyRequest{Name: "acme", Type: "gamma"})
	require.NoError(t, err)

	id, err := model.NewDummyID(dummyTestUUID)
	require.NoError(t, err)

	stored, err := repo.FindByID(context.Background(), *id)
	require.NoError(t, err)
	require.Equal(t, "acme", stored.Name().String())
	require.Equal(t, "gamma", stored.Type().String())
	require.Len(t, pub.events, 1)

	got, ok := pub.events[0].(*model.DummyCreated)
	require.True(t, ok)
	require.Equal(t, dummyTestUUID, got.AggregateID())
	require.Equal(t, "acme", got.DummyName())
	require.True(t, got.DummyCreatedAt().Equal(now))
}

func TestCreateDummyHandler_Handle_InvalidName_DoesNotPersistOrPublish(t *testing.T) {
	t.Parallel()

	repo := repository.NewMemoryDummyRepository()
	pub := &recordingPublisher{}
	h := command.NewCreateDummyHandler(
		repo,
		fixedUUIDGen{dummyTestUUID},
		fixedClock{instant: time.Now().UTC()},
		pub,
	)

	err := h.Handle(context.Background(), command.CreateDummyRequest{Name: "   ", Type: "alpha"})
	require.Error(t, err)
	require.ErrorIs(t, err, model.ErrDummyNameEmpty)
	require.Len(t, pub.events, 0)

	id := mustDummyID(t, dummyTestUUID)
	_, findErr := repo.FindByID(context.Background(), id)
	require.Error(t, findErr)
}

func TestCreateDummyHandler_Handle_InvalidType_DoesNotPersistOrPublish(t *testing.T) {
	t.Parallel()

	repo := repository.NewMemoryDummyRepository()
	pub := &recordingPublisher{}
	h := command.NewCreateDummyHandler(
		repo,
		fixedUUIDGen{dummyTestUUID},
		fixedClock{instant: time.Now().UTC()},
		pub,
	)

	err := h.Handle(context.Background(), command.CreateDummyRequest{Name: "ok", Type: "delta"})
	require.ErrorIs(t, err, model.ErrDummyTypeInvalid)
	require.Len(t, pub.events, 0)

	id := mustDummyID(t, dummyTestUUID)
	_, findErr := repo.FindByID(context.Background(), id)
	require.Error(t, findErr)
}

func TestCreateDummyHandler_Handle_RepoFails_DoesNotPublish(t *testing.T) {
	t.Parallel()

	pub := &recordingPublisher{}
	h := command.NewCreateDummyHandler(
		saveFailsRepo{},
		fixedUUIDGen{dummyTestUUID},
		fixedClock{instant: time.Now().UTC()},
		pub,
	)

	err := h.Handle(context.Background(), command.CreateDummyRequest{Name: "ok", Type: "beta"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "save failed")
	require.Empty(t, pub.events)
}

func mustDummyID(t *testing.T, s string) model.DummyID {
	t.Helper()
	id, err := model.NewDummyID(s)
	require.NoError(t, err)
	require.NotNil(t, id)
	return *id
}
