package model_test

import (
	"testing"

	"github.com/fermin/gophercraft/internal/domain/model"
)

func TestExample_IsValid(t *testing.T) {
	t.Parallel()

	e := model.Example{ID: "1", Name: "demo"}
	if !e.IsValid() {
		t.Fatalf("expected Example to be valid; got invalid for %#v", e)
	}
}
