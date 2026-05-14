package model_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fermin/gophercraft/internal/domain/model"
)

func TestDummy_IsValid(t *testing.T) {
	t.Parallel()

	e := model.Dummy{ID: "1", Name: "demo"}
	require.True(t, e.IsValid(), "expected Dummy to be valid for %#v", e)
}
