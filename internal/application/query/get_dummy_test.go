package query_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/fermin/gophercraft/internal/application/query"
	"github.com/fermin/gophercraft/internal/domain/model"
	"github.com/fermin/gophercraft/internal/infrastructure/repository"
)

const dummyTestUUID = "550e8400-e29b-41d4-a716-446655440000"

func TestGetDummyHandler_Handle_ReturnsStoredDummy(t *testing.T) {
	t.Parallel()

	repo := repository.NewMemoryDummyRepository()
	now := time.Date(2026, 5, 14, 15, 30, 0, 0, time.UTC)

	id := mustDummyID(t, dummyTestUUID)
	name := mustDummyName(t, "acme")
	dummyType := mustDummyType(t, "gamma")
	createdAt := mustDummyCreatedAt(t, now)

	d := model.CreateDummy(id, name, dummyType, createdAt)
	require.NoError(t, repo.Save(context.Background(), d))

	h := query.NewGetDummyHandler(repo)

	got, err := h.Handle(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, dummyTestUUID, got.ID().String())
	require.Equal(t, "acme", got.Name().String())
	require.Equal(t, "gamma", got.Type().String())
	require.True(t, got.CreatedAt().Time().Equal(now))
}

func TestGetDummyHandler_Handle_UnknownID_ReturnsError(t *testing.T) {
	t.Parallel()

	repo := repository.NewMemoryDummyRepository()
	h := query.NewGetDummyHandler(repo)

	id := mustDummyID(t, dummyTestUUID)

	_, err := h.Handle(context.Background(), id)
	require.Error(t, err)
	require.ErrorIs(t, err, model.ErrDummyNotFound)
}

func mustDummyID(t *testing.T, s string) model.DummyID {
	t.Helper()
	id, err := model.NewDummyID(s)
	require.NoError(t, err)
	require.NotNil(t, id)
	return *id
}

func mustDummyName(t *testing.T, s string) model.DummyName {
	t.Helper()
	name, err := model.NewDummyName(s)
	require.NoError(t, err)
	require.NotNil(t, name)
	return *name
}

func mustDummyType(t *testing.T, s string) model.DummyType {
	t.Helper()
	dummyType, err := model.NewDummyType(s)
	require.NoError(t, err)
	require.NotNil(t, dummyType)
	return *dummyType
}

func mustDummyCreatedAt(t *testing.T, instant time.Time) model.DummyCreatedAt {
	t.Helper()
	createdAt, err := model.NewDummyCreatedAt(instant)
	require.NoError(t, err)
	require.NotNil(t, createdAt)
	return *createdAt
}
