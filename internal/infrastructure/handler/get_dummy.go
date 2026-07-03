package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fermin/gophercraft/internal/domain/model"
)

type getDummyUseCase interface {
	Handle(ctx context.Context, id model.DummyID) (model.Dummy, error)
}

// GetDummyGinHandler adapts the get Dummy use case to Gin HTTP.
type GetDummyGinHandler struct {
	uc getDummyUseCase
}

// NewGetDummyGinHandler constructs a Gin adapter for the get Dummy use case.
func NewGetDummyGinHandler(uc getDummyUseCase) *GetDummyGinHandler {
	return &GetDummyGinHandler{uc: uc}
}

// Handle reads the path ID, invokes the use case, and maps errors to HTTP responses.
func (h *GetDummyGinHandler) Handle(c *gin.Context) {
	id, err := model.NewDummyID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dummy, err := h.uc.Handle(c.Request.Context(), *id)
	if err != nil {
		if errors.Is(err, model.ErrDummyIDEmpty) || errors.Is(err, model.ErrDummyIDInvalid) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, model.ErrDummyNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	p := dummy.ToPrimitives()
	c.JSON(http.StatusOK, gin.H{
		"id":         p.ID,
		"name":       p.Name,
		"type":       p.Type,
		"created_at": p.CreatedAt,
	})
}
