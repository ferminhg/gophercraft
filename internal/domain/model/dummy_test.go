package model_test

import (
	"testing"

	"github.com/fermin/gophercraft/internal/domain/model"
)

func TestDummy_IsValid(t *testing.T) {
	t.Parallel()

	e := model.Dummy{ID: "1", Name: "demo"}
	if !e.IsValid() {
		t.Fatalf("expected Dummy to be valid; got invalid for %#v", e)
	}
}
