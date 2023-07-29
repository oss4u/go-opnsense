package overrides

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/oss4u/go-opnsense/opnsense"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"testing"
)

type HostsTestSuite struct {
	suite.Suite
}

func (s HostsTestSuite) SetupTest() {
	workingDir, _ := os.Getwd()
	err := godotenv.Load(fmt.Sprintf("%s/.env", workingDir))

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func (s HostsTestSuite) TestCreateUpdateDelete() {
	if os.Getenv("OPNSENSE_ADDRESS") == "" {
		s.T().Skip("Missing credentials")
	}
	api := opnsense.GetOpnSenseClient("", "", "")
	overrides := GetHostsOverrideApi(api)
	hostDetails := OverridesHostDetails{
		Uuid:        "",
		Enabled:     true,
		Hostname:    "123",
		Domain:      "asdf",
		Rr:          "A",
		Description: "asdfasdf",
		Server:      "10.10.10.10",
	}
	host := OverridesHost{Host: hostDetails}
	createdHost, _ := overrides.Create(&host)
	createdHost.Host.Hostname = "456"
	updatedHost, _ := overrides.Update(createdHost)
	overrides.Delete(updatedHost)

}

func TestHostsTestSuite(t *testing.T) {
	suite.Run(t, new(HostsTestSuite))
}
