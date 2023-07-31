package overrides

import (
	"encoding/json"
	"github.com/oss4u/go-opnsense/opnsense/types"
	"github.com/stretchr/testify/assert"
)

func (suite OverridesTestSuite) TestFromJsonToOverridesAlias() {
	json_string := `{
	"alias": {
		"enabled": "1",
		"host": "fc3e8cfd-5ac1-4e7e-94d6-6e2b77ebeb2f",
		"hostname": "test",
		"domain": "int.sys-int.de",
		"description": "descr"
	}
}`
	alias := OverridesAlias{}
	err := json.Unmarshal([]byte(json_string), &alias)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), alias.Alias)
	assert.Equal(suite.T(), types.Bool(true), alias.Alias.Enabled)
	assert.Equal(suite.T(), "test", alias.Alias.Hostname)
	assert.Equal(suite.T(), "int.sys-int.de", alias.Alias.Domain)
	assert.Equal(suite.T(), "fc3e8cfd-5ac1-4e7e-94d6-6e2b77ebeb2f", alias.Alias.Host)
	assert.Equal(suite.T(), "descr", alias.Alias.Description)
}

func (suite OverridesTestSuite) TestFromJsonToOverridesAliasDetails() {
	json_string := `{
		"enabled": "1",
		"host": "fc3e8cfd-5ac1-4e7e-94d6-6e2b77ebeb2f",
		"hostname": "test",
		"domain": "int.sys-int.de",
		"description": "descr"
	}`
	alias := OverridesAliasDetails{}
	err := json.Unmarshal([]byte(json_string), &alias)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), types.Bool(true), alias.Enabled)
	assert.Equal(suite.T(), "test", alias.Hostname)
	assert.Equal(suite.T(), "int.sys-int.de", alias.Domain)
	assert.Equal(suite.T(), "fc3e8cfd-5ac1-4e7e-94d6-6e2b77ebeb2f", alias.Host)
	assert.Equal(suite.T(), "descr", alias.Description)
}
