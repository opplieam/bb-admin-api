package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/opplieam/bb-admin-api/internal/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type LoginUnitTestSuite struct {
	suite.Suite
}

func TestLoginHandler(t *testing.T) {
	suite.Run(t, new(LoginUnitTestSuite))
}

func (s *LoginUnitTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	utils.GetEnvForTesting()
}

func (s *LoginUnitTestSuite) TestLoginHandler() {
	testCases := []struct {
		name         string
		body         gin.H
		buildStubs   func(store *MockStorer)
		wantedText   string
		wantedStatus int
	}{
		{
			name: "valid credential",
			body: gin.H{"username": "admin", "password": "admin1234"},
			buildStubs: func(store *MockStorer) {
				store.EXPECT().FindByCredential(mock.Anything, mock.Anything).Return(2, nil).Once()
			},
			wantedText:   "token",
			wantedStatus: http.StatusOK,
		},
		{
			name: "wrong body request",
			body: gin.H{"user": "admin", "pass": "admin12345"},
			buildStubs: func(store *MockStorer) {
			},
			wantedText:   "wrong credentials",
			wantedStatus: http.StatusBadRequest,
		},
		{
			name: "wrong credential",
			body: gin.H{"username": "a", "password": "admin1234"},
			buildStubs: func(store *MockStorer) {
				store.EXPECT().
					FindByCredential(mock.Anything, mock.Anything).
					Return(0, fmt.Errorf("record not found")).
					Once()
			},
			wantedText:   "wrong credentials",
			wantedStatus: http.StatusNotFound,
		},
	}
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			mockDB := NewMockStorer(s.T())
			tc.buildStubs(mockDB)

			router := gin.Default()
			userH := NewHandler(mockDB)
			router.POST("/login", userH.LoginHandler)
			reqBody, _ := json.Marshal(tc.body)
			req, _ := http.NewRequest("POST", "/login", strings.NewReader(string(reqBody)))

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			s.Assert().Equal(w.Code, tc.wantedStatus)
			s.Assert().Contains(w.Body.String(), tc.wantedText)
		})
	}

}
