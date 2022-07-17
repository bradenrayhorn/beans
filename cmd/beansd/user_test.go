package main_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type RegisterSuite struct{ suite.Suite }

func TestRegisterSuite(t *testing.T) {
	suite.Run(t, new(RegisterSuite))
}

func (s *RegisterSuite) TestRegister() {
	ta := StartApplication(s.T())
	defer ta.Stop(s.T())

	r, err := ta.PostRequest("api/v1/user/register", map[string]interface{}{"username": "user", "password": "pass"})
	s.Require().Nil(err)
	s.Assert().Equal(http.StatusOK, r.StatusCode)
}
