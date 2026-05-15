package model_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/fermin/gophercraft/internal/domain/model"
)

func TestNewDummyCreated_ImplementsContract(t *testing.T) {
	t.Parallel()

	id := mustDummyID(t, testUUID)
	name := mustDummyName(t, "hello")
	createdAt := mustDummyCreatedAt(t, time.Date(2026, 5, 14, 12, 0, 0, 0, time.UTC))

	ev := model.NewDummyCreated(*id, *name, model.DummyTypeBeta, *createdAt)
	require.NotNil(t, ev)

	require.Equal(t, "dummy.created", ev.EventName())
	require.Equal(t, testUUID, ev.AggregateID())
	require.True(t, ev.OccurredAt().Equal(createdAt.Time()))

	require.Equal(t, testUUID, ev.DummyID())
	require.Equal(t, "hello", ev.DummyName())
	require.Equal(t, "beta", ev.DummyType())
	require.True(t, ev.DummyCreatedAt().Equal(createdAt.Time()))
}

func mustDummyID(t *testing.T, s string) *model.DummyID {
	t.Helper()
	id, err := model.NewDummyID(s)
	require.NoError(t, err)
	return id
}

func mustDummyName(t *testing.T, s string) *model.DummyName {
	t.Helper()
	n, err := model.NewDummyName(s)
	require.NoError(t, err)
	return n
}

func mustDummyCreatedAt(t *testing.T, tm time.Time) *model.DummyCreatedAt {
	t.Helper()
	c, err := model.NewDummyCreatedAt(tm)
	require.NoError(t, err)
	return c
}
