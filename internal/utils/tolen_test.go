package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type TokenSuite struct {
	suite.Suite
}

func TestToken(t *testing.T) {
	suite.Run(t, new(TokenSuite))
}

func (s *TokenSuite) SetupSuite() {
	GetEnvForTesting()
}

func (s *TokenSuite) TestValidToken() {
	token, err := GenerateToken(time.Hour, 2)
	s.Assert().NoError(err)
	s.Assert().NotEmpty(token)
	s.Assert().Contains(token, "v4.local")

	err = VerifyToken(token)
	s.Assert().NoError(err)
}

func (s *TokenSuite) TestInvalidToken() {
	token, err := GenerateToken(-1*time.Hour, 2)
	s.Assert().NoError(err)
	s.Assert().NotEmpty(token)
	s.Assert().Contains(token, "v4.local")

	err = VerifyToken(token)
	s.Assert().Error(err)
}
