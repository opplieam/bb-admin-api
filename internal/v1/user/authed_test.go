package user

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/opplieam/bb-admin-api/internal/store"
	"github.com/opplieam/bb-admin-api/internal/utils"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var runTest = flag.String(
	"run",
	"u",
	"u = only unit test, i = only integration test, ui = both unit and integration",
)

type AuthedUnitTestSuite struct {
	suite.Suite
}

func TestAuthenticateHandler(t *testing.T) {
	if strings.Contains(*runTest, "u") {
		suite.Run(t, new(AuthedUnitTestSuite))
	}
	if strings.Contains(*runTest, "i") {
		suite.Run(t, new(AuthedIntegrTestSuite))
	}
}

func (s *AuthedUnitTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	utils.GetEnvForTesting()
}

func (s *AuthedUnitTestSuite) TestLoginUnit() {
	testCases := []struct {
		name             string
		body             gin.H
		buildStubs       func(store *MockStorer)
		wantedText       string
		wantedStatus     int
		wantedSetCookies string
		wantedSameSite   string
	}{
		{
			name: "valid credential",
			body: gin.H{"username": "admin", "password": "admin1234"},
			buildStubs: func(store *MockStorer) {
				store.EXPECT().FindByCredential(mock.Anything, mock.Anything).Return(2, nil).Once()
			},
			wantedText:       "token",
			wantedStatus:     http.StatusOK,
			wantedSetCookies: "refresh_token",
			wantedSameSite:   "SameSite=Strict",
		},
		{
			name: "wrong body request",
			body: gin.H{"user": "admin", "pass": "admin12345"},
			buildStubs: func(store *MockStorer) {
			},
			wantedText:       "wrong credentials",
			wantedStatus:     http.StatusBadRequest,
			wantedSetCookies: "",
			wantedSameSite:   "",
		},
		{
			name: "wrong username",
			body: gin.H{"username": "a", "password": "admin1234"},
			buildStubs: func(store *MockStorer) {
				store.EXPECT().
					FindByCredential(mock.Anything, mock.Anything).
					Return(0, fmt.Errorf("record not found")).
					Once()
			},
			wantedText:       "wrong credentials",
			wantedStatus:     http.StatusNotFound,
			wantedSetCookies: "",
			wantedSameSite:   "",
		},
	}
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			mockStore := NewMockStorer(s.T())
			tc.buildStubs(mockStore)

			router := gin.Default()
			userH := NewHandler(mockStore)
			router.POST("/login", userH.LoginHandler)

			reqBody, err := json.Marshal(tc.body)
			s.Require().NoError(err)
			req, err := http.NewRequest("POST", "/login", bytes.NewReader(reqBody))
			s.Require().NoError(err)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			s.Assert().Equal(tc.wantedStatus, w.Code)
			s.Assert().Contains(w.Body.String(), tc.wantedText)
			s.Assert().Contains(w.Header().Get("Set-Cookie"), tc.wantedSetCookies)
		})
	}

}

func (s *AuthedUnitTestSuite) TestLogoutUnit() {
	mockStore := NewMockStorer(s.T())
	router := gin.Default()
	userH := NewHandler(mockStore)
	router.GET("/logout", userH.LogoutHandler)

	req := httptest.NewRequest("GET", "/logout", nil)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "v4.local.HN....."})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	setCookies := w.Header().Get("Set-Cookie")
	s.Assert().Contains(setCookies, "refresh_token=;")
	s.Assert().Contains(setCookies, "Max-Age=0;")
}

type AuthedIntegrTestSuite struct {
	suite.Suite
	TestDB     *sql.DB
	DockerPool *dockertest.Pool
	Resource   *dockertest.Resource
}

func (s *AuthedIntegrTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	utils.GetEnvForTesting()
}

func (s *AuthedIntegrTestSuite) SetupTest() {
	testDB, pool, resource, err := utils.CreateDockerTestContainer()
	s.Require().NoError(err, "failed to create container")

	// migrate database
	err = utils.MigrateDB(testDB, "file://../../../migrations/")
	s.Require().NoError(err, "failed to migrate database")

	// seed data
	err = utils.SeedData(testDB, "../../../data/test_user.sql")
	s.Require().NoError(err, "failed to seed data")

	s.DockerPool = pool
	s.Resource = resource
	s.TestDB = testDB
}

func (s *AuthedIntegrTestSuite) TestLoginIntegr() {
	testCases := []struct {
		name             string
		body             gin.H
		wantedText       string
		wantedStatus     int
		wantedSetCookies string
		wantedSameSite   string
	}{
		{
			name:             "valid credential",
			body:             gin.H{"username": "admin", "password": "admin1234"},
			wantedText:       "token",
			wantedStatus:     http.StatusOK,
			wantedSetCookies: "refresh_token",
			wantedSameSite:   "SameSite=Strict",
		},
		{
			name:             "wrong body request",
			body:             gin.H{"user": "admin", "pass": "admin12345"},
			wantedText:       "wrong credentials",
			wantedStatus:     http.StatusBadRequest,
			wantedSetCookies: "",
			wantedSameSite:   "",
		},
		{
			name:             "wrong password",
			body:             gin.H{"username": "admin", "password": "admin1111"},
			wantedText:       "wrong credentials",
			wantedStatus:     http.StatusNotFound,
			wantedSetCookies: "",
			wantedSameSite:   "",
		},
		{
			name:             "inactive user",
			body:             gin.H{"username": "pon", "password": "admin1234"},
			wantedText:       "wrong credentials",
			wantedStatus:     http.StatusNotFound,
			wantedSetCookies: "",
			wantedSameSite:   "",
		},
	}
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			router := gin.Default()
			userStore := store.NewUserStore(s.TestDB)
			userH := NewHandler(userStore)
			router.POST("/login", userH.LoginHandler)

			reqBody, err := json.Marshal(tc.body)
			s.Require().NoError(err)
			req, err := http.NewRequest("POST", "/login", bytes.NewReader(reqBody))
			s.Require().NoError(err)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			s.Assert().Equal(tc.wantedStatus, w.Code)
			s.Assert().Contains(w.Body.String(), tc.wantedText)
			s.Assert().Contains(w.Header().Get("Set-Cookie"), tc.wantedSetCookies)
		})
	}
}

func (s *AuthedIntegrTestSuite) TearDownTest() {
	err := s.DockerPool.Purge(s.Resource)
	s.Require().NoError(err, "could not purge pool")
}