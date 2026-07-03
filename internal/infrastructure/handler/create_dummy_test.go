package handler_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/fermin/gophercraft/internal/application/command"
	"github.com/fermin/gophercraft/internal/domain/model"
	"github.com/fermin/gophercraft/internal/infrastructure/handler"
)

type stubCreateDummyUseCase struct {
	lastReq command.CreateDummyRequest
	err     error
}

func (s *stubCreateDummyUseCase) Handle(_ context.Context, req command.CreateDummyRequest) error {
	s.lastReq = req
	return s.err
}

func TestCreateDummyGinHandler_Handle_Created(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	uc := &stubCreateDummyUseCase{}
	h := handler.NewCreateDummyGinHandler(uc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/dummies", bytes.NewBufferString(`{"name":"acme","type":"gamma"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Handle(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "acme", uc.lastReq.Name)
	assert.Equal(t, "gamma", uc.lastReq.Type)
}

func TestCreateDummyGinHandler_Handle_InvalidJSON_BadRequest(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	h := handler.NewCreateDummyGinHandler(&stubCreateDummyUseCase{})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/dummies", bytes.NewBufferString(`{invalid`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Handle(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateDummyGinHandler_Handle_UseCaseError_InternalServerError(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	uc := &stubCreateDummyUseCase{err: errors.New("save failed")}
	h := handler.NewCreateDummyGinHandler(uc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/dummies", bytes.NewBufferString(`{"name":"acme","type":"gamma"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Handle(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateDummyGinHandler_Handle_ValidationError_BadRequest(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	uc := &stubCreateDummyUseCase{err: model.ErrDummyNameEmpty}
	h := handler.NewCreateDummyGinHandler(uc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/dummies", bytes.NewBufferString(`{"name":"   ","type":"gamma"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	h.Handle(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
