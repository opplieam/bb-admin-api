package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Storer interface {
	IsAuthenticated(username, password string) error
}

type Handler struct {
	Store Storer
}

func NewHandler(store Storer) *Handler {
	return &Handler{
		Store: store,
	}
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
	err := h.Store.IsAuthenticated(loginParams.Username, loginParams.Password)
	if err != nil {
		_ = c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	c.JSON(200, gin.H{
		"username": loginParams.Username,
		"password": loginParams.Password,
		//"hashed":   string(hashPass),
	})
}
