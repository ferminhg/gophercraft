package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fermin/gophercraft/internal/application/command"
	"github.com/fermin/gophercraft/internal/domain/model"
)

type createDummyUseCase interface {
	Handle(ctx context.Context, req command.CreateDummyRequest) error
}

// CreateDummyGinHandler adapts the create Dummy use case to Gin HTTP.
type CreateDummyGinHandler struct {
	uc createDummyUseCase
}

// NewCreateDummyGinHandler constructs a Gin adapter for the create Dummy use case.
func NewCreateDummyGinHandler(uc createDummyUseCase) *CreateDummyGinHandler {
	return &CreateDummyGinHandler{uc: uc}
}

// Handle binds JSON, invokes the use case, and maps errors to HTTP responses.
func (h *CreateDummyGinHandler) Handle(c *gin.Context) {
	var body struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.uc.Handle(c.Request.Context(), command.CreateDummyRequest{
		Name: body.Name,
		Type: body.Type,
	})
	if err != nil {
		if isDomainValidationError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusCreated)
}

func isDomainValidationError(err error) bool {
	return errors.Is(err, model.ErrDummyIDEmpty) ||
		errors.Is(err, model.ErrDummyIDInvalid) ||
		errors.Is(err, model.ErrDummyNameEmpty) ||
		errors.Is(err, model.ErrDummyNameTooLong) ||
		errors.Is(err, model.ErrDummyTypeInvalid) ||
		errors.Is(err, model.ErrDummyCreatedAtZero)
}
