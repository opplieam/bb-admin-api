package middleware

import (
	"bytes"
	"log/slog"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	engine *gin.Engine
	buf    *bytes.Buffer
}

func TestSLogger(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) SetupTest() {
	buffer := new(bytes.Buffer)
	var memLogger = slog.New(slog.NewJSONHandler(buffer, nil))

	gin.SetMode(gin.TestMode)

	r := gin.New()

	r.Use(SLogger(memLogger))
	r.Use(gin.Recovery())

	r.GET("/valid", func(c *gin.Context) {})
	r.GET("/forceError", func(c *gin.Context) { panic("forced panic") })

	suite.engine = r
	suite.buf = buffer
}

func (suite *TestSuite) TestValid() {
	req := httptest.NewRequest("GET", "/valid?a=100", nil)
	w := httptest.NewRecorder()
	suite.engine.ServeHTTP(w, req)

	suite.Assert().Contains(suite.buf.String(), "200")
	suite.Assert().Contains(suite.buf.String(), "GET")
	suite.Assert().Contains(suite.buf.String(), "/valid")
	suite.Assert().Contains(suite.buf.String(), "a=100")

}

func (suite *TestSuite) TestServerError() {
	req := httptest.NewRequest("GET", "/forceError", nil)
	w := httptest.NewRecorder()
	suite.engine.ServeHTTP(w, req)

	suite.Assert().Contains(suite.buf.String(), "500")
	suite.Assert().Contains(suite.buf.String(), "GET")
	suite.Assert().Contains(suite.buf.String(), "/forceError")
	suite.Assert().Contains(suite.buf.String(), "ERROR")

}

func (suite *TestSuite) TestClientError() {
	req := httptest.NewRequest("GET", "/notfound", nil)
	w := httptest.NewRecorder()
	suite.engine.ServeHTTP(w, req)

	suite.Assert().Contains(suite.buf.String(), "404")
	suite.Assert().Contains(suite.buf.String(), "GET")
	suite.Assert().Contains(suite.buf.String(), "/notfound")
	suite.Assert().Contains(suite.buf.String(), "WARN")

}
