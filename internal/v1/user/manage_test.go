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

func (s *ManageUnitTestSuite) TestGetAllUsersUnit() {
	testData := []dbStore.AllUsersResult{
		{ID: 1, Username: faker.Username(), Active: true},
		{ID: 2, Username: faker.Username(), Active: false},
		{ID: 3, Username: faker.Username(), Active: true},
	}

	testCases := []struct {
		name             string
		buildStubs       func(store *MockStorer)
		wantedStatusCode int
		wantedBody       map[string][]dbStore.AllUsersResult
	}{
		{
			name: "successful get all users",
			buildStubs: func(store *MockStorer) {
				store.EXPECT().GetAllUsers().Return(testData, nil).Once()
			},
			wantedStatusCode: http.StatusOK,
			wantedBody: map[string][]dbStore.AllUsersResult{
				"data": testData,
			},
		},
		{
			name: "no record found",
			buildStubs: func(store *MockStorer) {
				store.EXPECT().GetAllUsers().Return(nil, dbStore.ErrRecordNotFound)
			},
			wantedStatusCode: http.StatusInternalServerError,
			wantedBody:       nil,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			mockStore := NewMockStorer(s.T())
			tc.buildStubs(mockStore)
			router := gin.Default()
			userH := NewHandler(mockStore)
			router.GET("/user", userH.GetAllUsers)

			req := httptest.NewRequest(http.MethodGet, "/user", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			var gotResponse map[string][]dbStore.AllUsersResult
			if tc.wantedBody != nil {
				err := json.Unmarshal(w.Body.Bytes(), &gotResponse)
				s.Require().NoError(err)
				s.Assert().Equal(tc.wantedBody, gotResponse)
			}

			s.Assert().Equal(tc.wantedStatusCode, w.Code)

		})
	}
}

func (s *ManageUnitTestSuite) TestUpdateUserStatusUnit() {
	testCases := []struct {
		name             string
		buildStubs       func(store *MockStorer)
		reqBody          gin.H
		wantedStatusCode int
	}{
		{
			name: "successful update",
			buildStubs: func(store *MockStorer) {
				store.EXPECT().UpdateUserStatus(mock.Anything, mock.Anything).Return(nil).Once()
			},
			reqBody:          gin.H{"id": 1, "active": true},
			wantedStatusCode: http.StatusNoContent,
		},
		{
			name:             "wrong request body",
			buildStubs:       func(store *MockStorer) {},
			reqBody:          gin.H{"user": "me"},
			wantedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "fail update",
			buildStubs: func(store *MockStorer) {
				store.EXPECT().UpdateUserStatus(mock.Anything, mock.Anything).Return(fmt.Errorf("some error")).Once()
			},
			reqBody:          gin.H{"id": 2, "active": false},
			wantedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			mockStore := NewMockStorer(s.T())
			tc.buildStubs(mockStore)
			router := gin.Default()
			userH := NewHandler(mockStore)
			router.PATCH("/user", userH.UpdateUserStatus)

			reqBody, err := json.Marshal(tc.reqBody)
			s.Require().NoError(err)
			req := httptest.NewRequest(http.MethodPatch, "/user", bytes.NewReader(reqBody))
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

func (s *ManageIntegrTestSuite) TestGetAllUsersIntegr() {
	testCases := []struct {
		name             string
		wantedStatusCode int
		wantedBody       []string
	}{
		{
			name:             "successful get all users",
			wantedStatusCode: http.StatusOK,
			wantedBody:       []string{"admin", "pon"},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			router := gin.Default()
			userH := NewHandler(dbStore.NewUserStore(s.TestDB))
			router.GET("/user", userH.GetAllUsers)

			req := httptest.NewRequest(http.MethodGet, "/user", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			s.Assert().Equal(tc.wantedStatusCode, w.Code)
			if tc.wantedBody != nil {
				for _, v := range tc.wantedBody {
					s.Assert().Contains(w.Body.String(), v)
				}
			}

		})
	}
}

func (s *ManageIntegrTestSuite) TearDownTest() {
	err := s.TestDB.Close()
	s.Require().NoError(err, "failed to close test database")
	err = s.DockerPool.Purge(s.Resource)
	s.Require().NoError(err, "could not purge pool")
}
