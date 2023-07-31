package overrides

import (
	"github.com/oss4u/go-opnsense/test"
	"github.com/stretchr/testify/suite"
	"testing"
)

type OverridesTestSuite struct {
	test.OpnsenseTestSuite
}

func TestHostsOverridesTestSuite(t *testing.T) {
	suite.Run(t, new(OverridesTestSuite))
}

func (s *OverridesTestSuite) SetupSuite() {
	s.SetupSuiteBase()
}

func (s *OverridesTestSuite) TearDownSuite() {
	s.TearDownSuiteBase()
}
