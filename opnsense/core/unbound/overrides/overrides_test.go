package overrides

import (
	"github.com/bmatcuk/go-vagrant"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
)

type OverridesTestSuite struct {
	suite.Suite
	vagrantClient *vagrant.VagrantClient
}

func TestHostsOverridesTestSuite(t *testing.T) {
	suite.Run(t, new(OverridesTestSuite))
}

func (s OverridesTestSuite) SetupSuite() {
	var err error
	err = godotenv.Load("../../../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	s.vagrantClient, err = vagrant.NewVagrantClient("../../../../")
	if err != nil {
		panic(err)
	}

	upcmd := s.vagrantClient.Up()
	upcmd.Verbose = true
	if err := upcmd.Run(); err != nil {
		panic(err)
	}
	if upcmd.Error != nil {
		panic(err)
	}

}

func (s OverridesTestSuite) TearDownSuite() {
	s.vagrantClient.Destroy()
}