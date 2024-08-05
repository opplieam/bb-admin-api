package user

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	dbStore "github.com/opplieam/bb-admin-api/internal/store"
	"github.com/opplieam/bb-admin-api/internal/utils"
	"github.com/ory/dockertest/v3"

	"github.com/gin-gonic/gin"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ManageUnitTestSuite struct {
	suite.Suite
}

func TestManageHandler(t *testing.T) {
	suite.Run(t, new(ManageUnitTestSuite))
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(ManageIntegrTestSuite))
}

func (s *ManageUnitTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (s *ManageUnitTestSuite) TestCreateUserUnit() {
	testCases := []struct {
		name             string
		buildStubs       func(store *MockStorer)
		reqBody          gin.H
		wantedStatusCode int
	}{
		{
			name: "successful creation",
			buildStubs: func(store *MockStorer) {
				store.EXPECT().CreateUser(mock.Anything, mock.Anything).Return(nil).Once()
			},
			reqBody:          gin.H{"username": faker.Username(), "password": faker.Password()},
			wantedStatusCode: http.StatusCreated,
		},
		{
			name:             "wrong request body",
			buildStubs:       func(store *MockStorer) {},
			reqBody:          gin.H{"user": faker.Username(), "pass": faker.Password()},
			wantedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "duplicate username",
			buildStubs: func(store *MockStorer) {
				store.EXPECT().CreateUser(mock.Anything, mock.Anything).Return(dbStore.ErrRecordAlreadyExists).Once()
			},
			reqBody:          gin.H{"username": faker.Username(), "password": faker.Password()},
			wantedStatusCode: http.StatusConflict,
		},
		{
			name: "internal server error",
			buildStubs: func(store *MockStorer) {
				store.EXPECT().CreateUser(mock.Anything, mock.Anything).Return(fmt.Errorf("other errors")).Once()
			},
			reqBody:          gin.H{"username": faker.Username(), "password": faker.Password()},
			wantedStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			mockStore := NewMockStorer(s.T())
			tc.buildStubs(mockStore)
			router := gin.Default()
			userH := NewHandler(mockStore)
			router.POST("/create", userH.CreateUser)

			reqBody, err := json.Marshal(tc.reqBody)
			s.Require().NoError(err)
			req, err := http.NewRequest(http.MethodPost, "/create", bytes.NewReader(reqBody))
			s.Require().NoError(err)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			s.Assert().Equal(tc.wantedStatusCode, w.Code)
		})
	}
}

// -----------------------------------------------------

type ManageIntegrTestSuite struct {
	suite.Suite
	TestDB     *sql.DB
	DockerPool *dockertest.Pool
	Resource   *dockertest.Resource
}

func (s *ManageIntegrTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
}

func (s *ManageIntegrTestSuite) SetupTest() {
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

func (s *ManageIntegrTestSuite) TestCreateUserIntegr() {
	testCases := []struct {
		name             string
		reqBody          gin.H
		wantedStatusCode int
	}{
		{
			name:             "successful creation",
			reqBody:          gin.H{"username": faker.Username(), "password": faker.Password()},
			wantedStatusCode: http.StatusCreated,
		},
		{
			name:             "wrong request body",
			reqBody:          gin.H{"user": faker.Username(), "pass": faker.Password()},
			wantedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:             "duplicate username",
			reqBody:          gin.H{"username": "admin", "password": faker.Password()},
			wantedStatusCode: http.StatusConflict,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			router := gin.Default()
			userH := NewHandler(dbStore.NewUserStore(s.TestDB))
			router.POST("/create", userH.CreateUser)

			reqBody, err := json.Marshal(tc.reqBody)
			s.Require().NoError(err)
			req, err := http.NewRequest(http.MethodPost, "/create", bytes.NewReader(reqBody))
			s.Require().NoError(err)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			s.Assert().Equal(tc.wantedStatusCode, w.Code)
		})
	}
}

func (s *ManageIntegrTestSuite) TearDownTest() {
	err := s.TestDB.Close()
	s.Require().NoError(err, "failed to close test database")
	err = s.DockerPool.Purge(s.Resource)
	s.Require().NoError(err, "could not purge pool")
}
