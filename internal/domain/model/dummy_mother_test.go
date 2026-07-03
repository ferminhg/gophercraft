package model_test

import (
	"testing"
	"time"

	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/require"

	"github.com/fermin/gophercraft/internal/domain/model"
	"github.com/fermin/gophercraft/internal/domain/port"
)

// DummyMother builds Dummy instances for tests using random fake data.
type DummyMother struct {
	tb      testing.TB
	uuidGen port.UUIDGenerator
}

// NewDummyMother constructs a DummyMother with the given UUID generator.
// Use FixedUUIDGenerator for deterministic tests or FakerUUIDGenerator for random ones.
func NewDummyMother(tb testing.TB, gen port.UUIDGenerator) DummyMother {
	tb.Helper()
	return DummyMother{tb: tb, uuidGen: gen}
}

// Random returns a valid Dummy with randomly generated field values.
func (m DummyMother) Random() model.Dummy {
	return m.build(func(b *dummyMotherBuild) {
		b.dummyType = faker.RandomElement(b.faker, model.DummyTypeAlpha, model.DummyTypeBeta, model.DummyTypeGamma)
	})
}

// WithType returns a valid Dummy with the given type and random other fields.
func (m DummyMother) WithType(t model.DummyType) model.Dummy {
	return m.build(func(b *dummyMotherBuild) {
		b.dummyType = t
	})
}

// WithName returns a valid Dummy with the given name and random other fields.
func (m DummyMother) WithName(name string) model.Dummy {
	return m.build(func(b *dummyMotherBuild) {
		b.nameStr = name
	})
}

type dummyMotherBuild struct {
	faker     faker.Faker
	nameStr   string
	dummyType model.DummyType
}

func (m DummyMother) build(customize func(*dummyMotherBuild)) model.Dummy {
	f := faker.New()
	b := &dummyMotherBuild{
		faker:     f,
		nameStr:   f.Person().Name(),
		dummyType: faker.RandomElement(f, model.DummyTypeAlpha, model.DummyTypeBeta, model.DummyTypeGamma),
	}
	customize(b)

	id, err := model.NewDummyID(m.uuidGen.Generate())
	require.NoError(m.tb, err)
	name, err := model.NewDummyName(b.nameStr)
	require.NoError(m.tb, err)
	createdAt, err := model.NewDummyCreatedAt(f.Time().Time(time.Now()))
	require.NoError(m.tb, err)
	return model.NewDummy(*id, *name, b.dummyType, *createdAt)
}
