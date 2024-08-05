package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opplieam/bb-admin-api/internal/store"
)

type ManageI interface {
	CreateUser(username, password string) error
}

func (h *Handler) CreateUser(c *gin.Context) {
	var loginParams loginParams
	if err := c.ShouldBindJSON(&loginParams); err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	err := h.Store.CreateUser(loginParams.Username, loginParams.Password)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrRecordAlreadyExists):
			_ = c.AbortWithError(http.StatusConflict, err)
			return
		default:
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	c.Status(http.StatusCreated)
}
