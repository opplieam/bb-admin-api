package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/opplieam/bb-admin-api/internal/utils"
)

type AuthHeader struct {
	Authorization string `header:"Authorization" binding:"required"`
}

func AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var header AuthHeader
		if err := c.ShouldBindHeader(&header); err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, err)
			c.JSON(-1, gin.H{"msg": "wrong header"})
			return
		}
		token := strings.Split(header.Authorization, "Bearer ")[1]

		err := utils.VerifyToken(token)
		if err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, err)
			c.JSON(-1, gin.H{"msg": errors.Unwrap(err).Error()})
			return
		}

		c.Next()
	}
}
