package overrides

import (
	"encoding/json"
	"github.com/oss4u/go-opnsense/opnsense/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromJsonToOverridesAlias(t *testing.T) {
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
	assert.Nil(t, err)
	assert.NotNil(t, alias.Alias)
	assert.Equal(t, types.Bool(true), alias.Alias.Enabled)
	assert.Equal(t, "test", alias.Alias.Hostname)
	assert.Equal(t, "int.sys-int.de", alias.Alias.Domain)
	assert.Equal(t, "fc3e8cfd-5ac1-4e7e-94d6-6e2b77ebeb2f", alias.Alias.Host)
	assert.Equal(t, "descr", alias.Alias.Description)
}

func TestFromJsonToOverridesAliasDetails(t *testing.T) {
	json_string := `{
		"enabled": "1",
		"host": "fc3e8cfd-5ac1-4e7e-94d6-6e2b77ebeb2f",
		"hostname": "test",
		"domain": "int.sys-int.de",
		"description": "descr"
	}`
	alias := OverridesAliasDetails{}
	err := json.Unmarshal([]byte(json_string), &alias)
	assert.Nil(t, err)
	assert.Equal(t, types.Bool(true), alias.Enabled)
	assert.Equal(t, "test", alias.Hostname)
	assert.Equal(t, "int.sys-int.de", alias.Domain)
	assert.Equal(t, "fc3e8cfd-5ac1-4e7e-94d6-6e2b77ebeb2f", alias.Host)
	assert.Equal(t, "descr", alias.Description)
}
