package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opplieam/bb-admin-api/internal/utils"
)

type LoginI interface {
	FindByCredential(username, password string) (int32, error)
}

type loginParams struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,gte=9"`
}

func (h *Handler) LoginHandler(c *gin.Context) {
	var loginParams loginParams
	if err := c.BindJSON(&loginParams); err != nil {
		c.JSON(-1, gin.H{"msg": "wrong credentials"})
		return
	}
	userId, err := h.Store.FindByCredential(loginParams.Username, loginParams.Password)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		c.JSON(-1, gin.H{"msg": "wrong credentials"})
		return
	}

	token, err := utils.GenerateToken(time.Hour, userId)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, gin.H{
		"token": token,
	})
}
