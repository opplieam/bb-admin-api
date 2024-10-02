package user

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opplieam/bb-admin-api/internal/store"
	"github.com/opplieam/bb-admin-api/internal/utils"
)

const (
	refreshTokenDuration = 730 * time.Hour // 1 month
	tokenDuration        = 1 * time.Hour
	cookiesAgeInt        = 2629800
)

type AuthedI interface {
	FindByCredential(username, password string) (int32, error)
	IsValidUser(id int32) error
}

type loginReqBody struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,gte=9"`
}

func (h *Handler) LoginHandler(c *gin.Context) {
	var loginRB loginReqBody
	if err := c.ShouldBindJSON(&loginRB); err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	userId, err := h.Store.FindByCredential(loginRB.Username, loginRB.Password)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	token, err := utils.GenerateToken(tokenDuration, userId)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	refreshToken, err := utils.GenerateToken(refreshTokenDuration, userId) // 1 month
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	//TODO: change domain name and secure depend on environment
	c.SetCookie(
		"refresh_token",
		refreshToken,
		cookiesAgeInt,
		"/",
		"localhost",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (h *Handler) LogoutHandler(c *gin.Context) {
	//TODO: change domain name and secure depend on environment
	c.SetCookie(
		"refresh_token",
		"", -1,
		"/", "localhost",
		false,
		true,
	)
	c.JSON(http.StatusOK, gin.H{"msg": "logged out"})
}

func (h *Handler) RefreshTokenHandler(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		c.JSON(-1, gin.H{"msg": "no token"})
		return
	}
	token, err := utils.VerifyToken(refreshToken)
	if err != nil {
		_ = c.AbortWithError(http.StatusUnauthorized, err)
		c.JSON(-1, gin.H{"msg": "invalid token"})
		return
	}

	userIdString, err := token.GetString("user_id")
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		c.JSON(-1, gin.H{"msg": "no user id"})
		return
	}
	userId, err := strconv.ParseInt(userIdString, 10, 32)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	err = h.Store.IsValidUser(int32(userId))
	if err != nil {
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			c.AbortWithStatus(http.StatusForbidden)
			return
		default:
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	newToken, err := utils.GenerateToken(tokenDuration, int32(userId))
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": newToken})

}
