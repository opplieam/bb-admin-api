package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

type loginParams struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,gte=9"`
}

func (h *Handler) LoginHandler(c *gin.Context) {
	var loginParams loginParams
	if err := c.BindJSON(&loginParams); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "wrong credentials"})
		return
	}

	c.JSON(200, loginParams)
}
