package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fermin/gophercraft/internal/domain/model"
	"github.com/fermin/gophercraft/internal/infrastructure/handler"
)

const dummyTestUUID = "550e8400-e29b-41d4-a716-446655440000"

type stubGetDummyUseCase struct {
	lastID model.DummyID
	dummy  model.Dummy
	err    error
}

func (s *stubGetDummyUseCase) Handle(_ context.Context, id model.DummyID) (model.Dummy, error) {
	s.lastID = id
	if s.err != nil {
		return model.Dummy{}, s.err
	}
	return s.dummy, nil
}

func TestGetDummyGinHandler_Handle_OK(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	now := time.Date(2026, 5, 14, 15, 30, 0, 0, time.UTC)
	id := mustDummyID(t, dummyTestUUID)
	name := mustDummyName(t, "acme")
	dummyType := mustDummyType(t, "gamma")
	createdAt := mustDummyCreatedAt(t, now)
	dummy := model.NewDummy(id, name, dummyType, createdAt)

	uc := &stubGetDummyUseCase{dummy: dummy}
	h := handler.NewGetDummyGinHandler(uc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: dummyTestUUID}}
	c.Request = httptest.NewRequest(http.MethodGet, "/dummies/"+dummyTestUUID, nil)

	h.Handle(c)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, dummyTestUUID, uc.lastID.String())

	var body struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		Type      string    `json:"type"`
		CreatedAt time.Time `json:"created_at"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &body)
	require.NoError(t, err)
	assert.Equal(t, dummyTestUUID, body.ID)
	assert.Equal(t, "acme", body.Name)
	assert.Equal(t, "gamma", body.Type)
	assert.True(t, body.CreatedAt.Equal(now))
}

func TestGetDummyGinHandler_Handle_InvalidID_BadRequest(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	h := handler.NewGetDummyGinHandler(&stubGetDummyUseCase{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "not-a-uuid"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/dummies/not-a-uuid", nil)

	h.Handle(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetDummyGinHandler_Handle_NotFound(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	uc := &stubGetDummyUseCase{err: model.ErrDummyNotFound}
	h := handler.NewGetDummyGinHandler(uc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: dummyTestUUID}}
	c.Request = httptest.NewRequest(http.MethodGet, "/dummies/"+dummyTestUUID, nil)

	h.Handle(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
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
