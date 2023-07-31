package test

import (
	"github.com/bmatcuk/go-vagrant"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
)

type OpnsenseTestSuite struct {
	suite.Suite
	vagrantClient   *vagrant.VagrantClient
	IntegrationTest bool
}

func (s *OpnsenseTestSuite) SetupSuiteBase() {
	var err error
	err = godotenv.Load("../../../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	s.IntegrationTest = false
	if os.Getenv("GITHUB_SHA") == "" {
		s.IntegrationTest = true
	}
	if s.IntegrationTest {
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

}

func (s *OpnsenseTestSuite) TearDownSuiteBase() {
	if s.IntegrationTest {
		s.vagrantClient.Destroy()
	}
}
