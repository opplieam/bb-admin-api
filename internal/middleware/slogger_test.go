package middleware

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type SLoggerTestSuite struct {
	suite.Suite
	engine *gin.Engine
	buf    *bytes.Buffer
}

func TestSLogger(t *testing.T) {
	suite.Run(t, new(SLoggerTestSuite))
}

func (s *SLoggerTestSuite) SetupTest() {
	buffer := new(bytes.Buffer)
	var memLogger = slog.New(slog.NewJSONHandler(buffer, nil))

	gin.SetMode(gin.TestMode)

	r := gin.New()

	r.Use(SLogger(memLogger))
	r.Use(gin.Recovery())

	s.engine = r
	s.buf = buffer
}

func (s *SLoggerTestSuite) isContain(args ...string) {
	for _, word := range args {
		s.Assert().Contains(s.buf.String(), word)
	}
}

func (s *SLoggerTestSuite) TestValid() {
	s.engine.GET("/valid", func(c *gin.Context) {})

	req := httptest.NewRequest("GET", "/valid?a=100", nil)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	s.isContain("200", "GET", "/valid", "a=100")
}

func (s *SLoggerTestSuite) TestServerError() {
	s.engine.GET("/forceError", func(c *gin.Context) { panic("forced panic") })

	req := httptest.NewRequest("GET", "/forceError", nil)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	s.isContain("500", "GET", "/forceError", "ERROR")
}

func (s *SLoggerTestSuite) TestClientError() {
	req := httptest.NewRequest("GET", "/notfound", nil)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	s.isContain("404", "GET", "/notfound", "WARN")
}

func (s *SLoggerTestSuite) TestAbortClientError() {
	s.engine.GET("/clientError", func(c *gin.Context) {
		_ = c.AbortWithError(http.StatusUnauthorized, errors.New("unauthorized"))
		return
	})

	req := httptest.NewRequest("GET", "/clientError", nil)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	s.isContain("client error", "unauthorized", "WARN")
}

func (s *SLoggerTestSuite) TestAbortServerError() {
	s.engine.GET("/serverError", func(c *gin.Context) {
		_ = c.AbortWithError(http.StatusInternalServerError, errors.New("internal error"))
	})
	req := httptest.NewRequest("GET", "/serverError", nil)
	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	s.isContain("server error", "internal error", "ERROR")
}
