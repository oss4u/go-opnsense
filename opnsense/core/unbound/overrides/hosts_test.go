package overrides

import (
	"github.com/oss4u/go-opnsense/opnsense"
	"github.com/stretchr/testify/assert"
)

func (s OverridesTestSuite) TestCreateUpdateDelete() {
	if !s.integrationTest {
		s.T().Skip("CI Build - Skipping integration tests")
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
	assert.NotNil(s.T(), createdHost)
	createdHost.Host.Hostname = "456"
	updatedHost, _ := overrides.Update(createdHost)
	assert.NotNil(s.T(), updatedHost)
	overrides.Delete(updatedHost)
	unavailableHost, _ := overrides.Read(updatedHost.Host.Uuid)
	assert.Nil(s.T(), unavailableHost)

}
