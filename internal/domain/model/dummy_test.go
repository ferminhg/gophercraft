package model_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/fermin/gophercraft/internal/domain/model"
)

const testUUID = "550e8400-e29b-41d4-a716-446655440000"

func TestDummyMother_Random(t *testing.T) {
	t.Parallel()

	d := NewDummyMother(NewFakerUUIDGenerator()).Random()
	require.NotEmpty(t, d.ID().String())
	require.NotEmpty(t, d.Name().String())
	require.True(t, d.Type().IsValid())
	require.False(t, d.CreatedAt().Time().IsZero())
}

func TestDummyMother_WithType(t *testing.T) {
	t.Parallel()

	d := NewDummyMother(NewFakerUUIDGenerator()).WithType(model.DummyTypeBeta)
	require.Equal(t, model.DummyTypeBeta, d.Type())
}

func TestDummyMother_WithName(t *testing.T) {
	t.Parallel()

	d := NewDummyMother(NewFakerUUIDGenerator()).WithName("fixed-name")
	require.Equal(t, "fixed-name", d.Name().String())
}

func TestDummyMother_FixedID(t *testing.T) {
	t.Parallel()

	d := NewDummyMother(FixedUUIDGenerator{Value: testUUID}).Random()
	require.Equal(t, testUUID, d.ID().String())
}

func TestNewDummy_Valid(t *testing.T) {
	t.Parallel()

	id, err := model.NewDummyID(testUUID)
	require.NoError(t, err)
	name, err := model.NewDummyName("hello")
	require.NoError(t, err)
	created, err := model.NewDummyCreatedAt(time.Now().UTC())
	require.NoError(t, err)

	d := model.NewDummy(id, name, model.DummyTypeAlpha, created)
	require.Equal(t, id, d.ID())
	require.Equal(t, name, d.Name())
	require.Equal(t, model.DummyTypeAlpha, d.Type())
	require.Equal(t, created, d.CreatedAt())
}

func TestNewDummyType_Valid(t *testing.T) {
	t.Parallel()

	for _, raw := range []string{"alpha", " beta ", "gamma"} {
		got, err := model.NewDummyType(raw)
		require.NoError(t, err)
		require.True(t, got.IsValid())
		require.NotEmpty(t, got.String())
	}
}

func TestNewDummyType_Invalid(t *testing.T) {
	t.Parallel()

	for _, raw := range []string{"", "unknown", "delta"} {
		_, err := model.NewDummyType(raw)
		require.Error(t, err)
	}
}

func TestDummyType_IsValid(t *testing.T) {
	t.Parallel()

	require.True(t, model.DummyTypeAlpha.IsValid())
	require.True(t, model.DummyTypeBeta.IsValid())
	require.True(t, model.DummyTypeGamma.IsValid())
	require.False(t, model.DummyType{}.IsValid())
}

func TestDummy_ToPrimitives(t *testing.T) {
	t.Parallel()

	id, err := model.NewDummyID(testUUID)
	require.NoError(t, err)
	name, err := model.NewDummyName("hello")
	require.NoError(t, err)
	created, err := model.NewDummyCreatedAt(time.Date(2025, 3, 14, 12, 0, 0, 0, time.UTC))
	require.NoError(t, err)

	d := model.NewDummy(id, name, model.DummyTypeBeta, created)
	prim := d.ToPrimitives()

	require.Equal(t, testUUID, prim.ID)
	require.Equal(t, "hello", prim.Name)
	require.Equal(t, "beta", prim.Type)
	require.True(t, prim.CreatedAt.Equal(created.Time()))
}

func TestDummyFromPrimitives_Valid(t *testing.T) {
	t.Parallel()

	created := time.Date(2025, 5, 1, 9, 30, 0, 0, time.UTC)
	d, err := model.DummyFromPrimitives(model.DummyPrimitives{
		ID:        testUUID,
		Name:      "aggregate-name",
		Type:      "gamma",
		CreatedAt: created,
	})
	require.NoError(t, err)
	require.Equal(t, testUUID, d.ID().String())
	require.Equal(t, "aggregate-name", d.Name().String())
	require.Equal(t, "gamma", d.Type().String())
	require.True(t, d.CreatedAt().Time().Equal(created))
}

func TestDummyFromPrimitives_Invalid(t *testing.T) {
	t.Parallel()

	base := model.DummyPrimitives{
		ID:        testUUID,
		Name:      "ok",
		Type:      "alpha",
		CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	tests := []struct {
		name string
		p    model.DummyPrimitives
	}{{
		name: "empty id",
		p: model.DummyPrimitives{
			ID:        "",
			Name:      base.Name,
			Type:      base.Type,
			CreatedAt: base.CreatedAt,
		},
	}, {
		name: "empty name",
		p: model.DummyPrimitives{
			ID:        base.ID,
			Name:      "",
			Type:      base.Type,
			CreatedAt: base.CreatedAt,
		},
	}, {
		name: "invalid type",
		p: model.DummyPrimitives{
			ID:        base.ID,
			Name:      base.Name,
			Type:      "delta",
			CreatedAt: base.CreatedAt,
		},
	}, {
		name: "zero created at",
		p: model.DummyPrimitives{
			ID:        base.ID,
			Name:      base.Name,
			Type:      base.Type,
			CreatedAt: time.Time{},
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := model.DummyFromPrimitives(tt.p)
			require.Error(t, err)
		})
	}
}
