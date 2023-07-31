package overrides

import (
	"encoding/json"
	"github.com/kinbiko/jsonassert"
	"github.com/oss4u/go-opnsense/opnsense/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HostsOverridesTestSuite struct {
	suite.Suite
}

func (s OverridesTestSuite) TestToJson() {
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
	//json := host.Host.ConvertToJson()
	data, err := json.Marshal(host)
	assert.Nil(s.T(), err)
	ja := jsonassert.New(s.T())
	ja.Assertf(string(data), `
	{
		"host": 
			{
				"enabled": "1",
				"hostname": "123",
				"rr": "A",
				"server": "10.10.10.10",
				"description": "asdfasdf",
				"domain": "asdf"
			}
	}`)
	//fmt.Print(string(result))

}

func (s OverridesTestSuite) TestFromJsonToOverridesHost() {
	jsonString := "{\"host\":{\"enabled\":\"1\",\"hostname\":\"srv01\",\"domain\":\"dev.sys-int.de\",\"rr\":{\"A\":{\"value\":\"A (IPv4 address)\",\"selected\":0},\"AAAA\":{\"value\":\"AAAA (IPv6 address)\",\"selected\":0},\"MX\":{\"value\":\"MX (Mail server)\",\"selected\":1}},\"mxprio\":\"10\",\"mx\":\"srv01.dev.sys-int.de\",\"server\":\"server01\",\"description\":\"srv01 - MX\"}}"
	host := OverridesHost{}
	err := json.Unmarshal([]byte(jsonString), &host)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), host.Host)
	assert.Equal(s.T(), types.Bool(true), host.Host.Enabled)
	assert.Equal(s.T(), "srv01", host.Host.Hostname)
	assert.Equal(s.T(), "dev.sys-int.de", host.Host.Domain)
	assert.Equal(s.T(), Rr("MX"), host.Host.Rr)
	assert.Equal(s.T(), "srv01.dev.sys-int.de", host.Host.Mx)
	assert.Equal(s.T(), MxPrio(10), host.Host.Mxprio)
	assert.Equal(s.T(), "srv01 - MX", host.Host.Description)
	assert.Equal(s.T(), "server01", host.Host.Server)
}

func (s OverridesTestSuite) TestFromJsonToOverridesHostDetails() {
	jsonString := "{\"enabled\":\"1\",\"hostname\":\"srv01\",\"domain\":\"dev.sys-int.de\",\"rr\":{\"A\":{\"value\":\"A (IPv4 address)\",\"selected\":0},\"AAAA\":{\"value\":\"AAAA (IPv6 address)\",\"selected\":0},\"MX\":{\"value\":\"MX (Mail server)\",\"selected\":1}},\"mxprio\":\"10\",\"mx\":\"srv01.dev.sys-int.de\",\"server\":\"server01\",\"description\":\"srv01 - MX\"}"
	host := OverridesHostDetails{}
	err := json.Unmarshal([]byte(jsonString), &host)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), types.Bool(true), host.Enabled)
	assert.Equal(s.T(), "srv01", host.Hostname)
	assert.Equal(s.T(), "dev.sys-int.de", host.Domain)
	assert.Equal(s.T(), Rr("MX"), host.Rr)
	assert.Equal(s.T(), "srv01.dev.sys-int.de", host.Mx)
	assert.Equal(s.T(), MxPrio(10), host.Mxprio)
	assert.Equal(s.T(), "srv01 - MX", host.Description)
	assert.Equal(s.T(), "server01", host.Server)
}

func (s OverridesTestSuite) TestMxPrioToInt() {
	var cut MxPrio
	cut = 101
	assert.Equal(s.T(), 101, cut.Int())
}

func (s OverridesTestSuite) TestRrToString() {
	var cut Rr
	cut = "AAAA"
	assert.Equal(s.T(), "AAAA", cut.String())
}
