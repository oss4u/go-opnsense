package overrides

import (
	"github.com/oss4u/go-opnsense/opnsense"
	"github.com/stretchr/testify/assert"
)

func (s OverridesTestSuite) TestCreateUpdateDeleteAliases() {
	if !s.IntegrationTest {
		s.T().Skip("CI Build - Skipping integration tests")
	}
	api := opnsense.GetOpnSenseClient("", "", "")
	hostsOverrideApi := GetHostsOverrideApi(api)
	aliasOverrideApi := GetAliasesOverrideApi(api)
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
	createdHost, _ := hostsOverrideApi.Create(&host)
	assert.NotNil(s.T(), createdHost)
	alias := OverridesAlias{Alias: OverridesAliasDetails{
		Enabled:     true,
		Host:        createdHost.Host.Uuid,
		Hostname:    "abc",
		Domain:      "asdddd",
		Description: "Alias",
	}}
	createdAlias, _ := aliasOverrideApi.Create(&alias)
	createdAlias.Alias.Hostname = "456"
	updatedAlias, _ := aliasOverrideApi.Update(createdAlias)
	assert.NotNil(s.T(), updatedAlias)
	aliasOverrideApi.Delete(updatedAlias)
	unavailableHost, _ := aliasOverrideApi.Read(updatedAlias.Alias.Uuid)
	assert.Nil(s.T(), unavailableHost)

}
