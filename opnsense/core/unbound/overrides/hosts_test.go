package overrides

import (
	"github.com/oss4u/go-opnsense/opnsense"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type HostsTestSuite struct {
	suite.Suite
}

func TestCreateUpdateDelete(t *testing.T) {
	if os.Getenv("OPNSENSE_ADDRESS") == "" {
		t.Skip("Missing credentials")
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
