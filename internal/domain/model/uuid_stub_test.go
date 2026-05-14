package model_test

import "github.com/jaswdr/faker/v2"

// FixedUUIDGenerator always returns the same UUID value.
// Use it in tests that need deterministic, assertable identifiers.
type FixedUUIDGenerator struct{ Value string }

func (f FixedUUIDGenerator) Generate() string { return f.Value }

// FakerUUIDGenerator generates a random UUID v4 via the faker library.
// Use it when the exact value does not matter but validity is required.
type FakerUUIDGenerator struct{ f faker.Faker }

func NewFakerUUIDGenerator() FakerUUIDGenerator {
	return FakerUUIDGenerator{f: faker.New()}
}

func (g FakerUUIDGenerator) Generate() string { return g.f.UUID().V4() }
