package user

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opplieam/bb-admin-api/internal/utils"
)

type AuthedI interface {
	FindByCredential(username, password string) (int32, error)
}

type loginParams struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,gte=9"`
}

func (h *Handler) LoginHandler(c *gin.Context) {
	var loginParams loginParams
	if err := c.ShouldBindJSON(&loginParams); err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	userId, err := h.Store.FindByCredential(loginParams.Username, loginParams.Password)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	token, err := utils.GenerateToken(time.Hour, userId)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	refreshToken, err := utils.GenerateToken(time.Hour*730, userId) // 1 month
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	//TODO: change domain name and secure depend on environment
	c.SetCookie(
		"refresh_token",
		refreshToken,
		2629800,
		"/",
		"localhost",
		false,
		true,
	)

	c.JSON(200, gin.H{
		"token": token,
	})
}

func (h *Handler) LogoutHandler(c *gin.Context) {
	//TODO: change domain name and secure depend on environment
	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{"msg": "logged out"})
}
