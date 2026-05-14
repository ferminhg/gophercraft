// Package model defines domain entities and invariants.
package model

import "time"

// DummyPrimitives carries primitive values for serializing or reconstructing a Dummy.
type DummyPrimitives struct {
	ID        string
	Name      string
	Type      string
	CreatedAt time.Time
}

// Dummy is a placeholder aggregate root for bootstrapping the layout.
type Dummy struct {
	id        DummyID
	name      DummyName
	dummyType DummyType
	createdAt DummyCreatedAt
}

// NewDummy assembles a Dummy from validated value objects.
func NewDummy(id DummyID, name DummyName, t DummyType, createdAt DummyCreatedAt) Dummy {
	return Dummy{
		id:        id,
		name:      name,
		dummyType: t,
		createdAt: createdAt,
	}
}

// ID returns the aggregate identifier.
func (d Dummy) ID() DummyID {
	return d.id
}

// Name returns the human-readable name.
func (d Dummy) Name() DummyName {
	return d.name
}

// Type returns the dummy classification.
func (d Dummy) Type() DummyType {
	return d.dummyType
}

// CreatedAt returns creation timestamp.
func (d Dummy) CreatedAt() DummyCreatedAt {
	return d.createdAt
}

// DummyFromPrimitives validates primitives and builds a Dummy.
func DummyFromPrimitives(p DummyPrimitives) (Dummy, error) {
	id, err := NewDummyID(p.ID)
	if err != nil {
		return Dummy{}, err
	}
	name, err := NewDummyName(p.Name)
	if err != nil {
		return Dummy{}, err
	}
	t, err := NewDummyType(p.Type)
	if err != nil {
		return Dummy{}, err
	}
	createdAt, err := NewDummyCreatedAt(p.CreatedAt)
	if err != nil {
		return Dummy{}, err
	}
	return NewDummy(id, name, t, createdAt), nil
}

// ToPrimitives returns the aggregate fields as plain values.
func (d Dummy) ToPrimitives() DummyPrimitives {
	return DummyPrimitives{
		ID:        d.id.String(),
		Name:      d.name.String(),
		Type:      d.dummyType.String(),
		CreatedAt: d.createdAt.Time(),
	}
}
