package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opplieam/bb-admin-api/internal/store"
)

type ManageI interface {
	CreateUser(username, password string) error
	GetAllUsers() ([]store.AllUsersResult, error)
	UpdateUserStatus(userId int32, active bool) error
}

func (h *Handler) CreateUser(c *gin.Context) {
	var loginRB loginReqBody
	if err := c.ShouldBindJSON(&loginRB); err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	err := h.Store.CreateUser(loginRB.Username, loginRB.Password)
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

func (h *Handler) GetAllUsers(c *gin.Context) {
	result, err := h.Store.GetAllUsers()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

type updateUserReqBody struct {
	ID     int32 `json:"id" binding:"required"`
	Active *bool `json:"active" binding:"required"`
}

func (h *Handler) UpdateUserStatus(c *gin.Context) {
	var updateUserRB updateUserReqBody
	if err := c.ShouldBindJSON(&updateUserRB); err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	err := h.Store.UpdateUserStatus(updateUserRB.ID, *updateUserRB.Active)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}
