package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/opplieam/bb-admin-api/internal/utils"
	"github.com/stretchr/testify/suite"
)

var validToken string = "v4.local.0qxB-NQ2FwO9EORIvbSDXGGOPBRM8pp20WvgZW4tfChl14XYrKur148KkpNm5aDrEIV-PQOSEjJ3iQjrXWTIXgik36a5vgKu0-XmJIbiowAm8vOoC6j3QeAV7YVtayN-rHO9tbAmN0b-MSiLxtwlPxhClaqqxUIfHBwanh2TLIIzU33ZuhS1eOOiNZd4-gaAizAX9FcvRhLi-rd_jkdsMwvCj6I"
var expireToken string = "v4.local.SCRzrl5yVNpxYC2ndTPmnhcdHqLhLB3Bw_ImOZOc0kMVrP_WLPJAjLF3-YSX0MepdhK1qnROonokFLlQBRQy2XB8M9gZZRpfZE8nQJ8TZCSGxvRZD6KV5-awMmEvqgCEfOCBAMWupAJd9ohkZWZZ-m-D_F5eMkFlZrge0wIJQjf2pTl5MV-NZUfNG9C-Kq6ImnASkX2_AZYk7QRknPYe0vFDJZQ"

type AuthorizationTestSuite struct {
	suite.Suite
}

func TestAuthorizationMiddleware(t *testing.T) {
	suite.Run(t, new(AuthorizationTestSuite))
}

func (s *AuthorizationTestSuite) SetupSuite() {
	utils.GetEnvForTesting()
	gin.SetMode(gin.TestMode)
}

func (s *AuthorizationTestSuite) TestValid() {
	r := gin.New()
	r.Use(AuthorizationMiddleware())
	r.POST("/valid", func(c *gin.Context) {})

	req := httptest.NewRequest("POST", "/valid", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", validToken))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusOK, w.Code)
}

func (s *AuthorizationTestSuite) TestInvalid() {
	testCases := []struct {
		name      string
		authValue string
		expCode   int
	}{
		{
			name:      "expired token",
			authValue: fmt.Sprintf("Bearer %s", expireToken),
			expCode:   http.StatusUnauthorized,
		},
		{
			name:      "invalid token",
			authValue: fmt.Sprintf("Bearer %s", expireToken+"not-a-real-token"),
			expCode:   http.StatusUnauthorized,
		},
		{
			name:      "not bearer token",
			authValue: fmt.Sprintf("JWT %s", validToken),
			expCode:   http.StatusUnauthorized,
		},
		{
			name:      "empty auth value",
			authValue: "",
			expCode:   http.StatusUnauthorized,
		},
		{
			name:      "no Authorization header",
			authValue: "no header",
			expCode:   http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			r := gin.New()
			r.Use(AuthorizationMiddleware())
			r.POST("/invalid", func(c *gin.Context) {})

			req := httptest.NewRequest("POST", "/invalid", nil)
			if tc.authValue != "no header" {
				req.Header.Set("Authorization", tc.authValue)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			s.Assert().Equal(tc.expCode, w.Code)
		})
	}

}
