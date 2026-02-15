package overrides

import (
	"encoding/json"
	"testing"

	"github.com/kinbiko/jsonassert"
	"github.com/oss4u/go-opnsense/opnsense/types"
	"github.com/stretchr/testify/assert"
)

func TestToJson(t *testing.T) {
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
	assert.Nil(t, err)
	ja := jsonassert.New(t)
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

func TestFromJsonToOverridesHost(t *testing.T) {
	jsonString := "{\"host\":{\"enabled\":\"1\",\"hostname\":\"srv01\",\"domain\":\"dev.sys-int.de\",\"rr\":{\"A\":{\"value\":\"A (IPv4 address)\",\"selected\":0},\"AAAA\":{\"value\":\"AAAA (IPv6 address)\",\"selected\":0},\"MX\":{\"value\":\"MX (Mail server)\",\"selected\":1}},\"mxprio\":\"10\",\"mx\":\"srv01.dev.sys-int.de\",\"server\":\"server01\",\"description\":\"srv01 - MX\"}}"
	host := OverridesHost{}
	err := json.Unmarshal([]byte(jsonString), &host)
	assert.Nil(t, err)
	assert.NotNil(t, host.Host)
	assert.Equal(t, types.Bool(true), host.Host.Enabled)
	assert.Equal(t, "srv01", host.Host.Hostname)
	assert.Equal(t, "dev.sys-int.de", host.Host.Domain)
	assert.Equal(t, Rr("MX"), host.Host.Rr)
	assert.Equal(t, "srv01.dev.sys-int.de", host.Host.Mx)
	assert.Equal(t, MxPrio(10), host.Host.Mxprio)
	assert.Equal(t, "srv01 - MX", host.Host.Description)
	assert.Equal(t, "server01", host.Host.Server)
}

func TestFromJsonToOverridesHostDetails(t *testing.T) {
	jsonString := "{\"enabled\":\"1\",\"hostname\":\"srv01\",\"domain\":\"dev.sys-int.de\",\"rr\":{\"A\":{\"value\":\"A (IPv4 address)\",\"selected\":0},\"AAAA\":{\"value\":\"AAAA (IPv6 address)\",\"selected\":0},\"MX\":{\"value\":\"MX (Mail server)\",\"selected\":1}},\"mxprio\":\"10\",\"mx\":\"srv01.dev.sys-int.de\",\"server\":\"server01\",\"description\":\"srv01 - MX\"}"
	host := OverridesHostDetails{}
	err := json.Unmarshal([]byte(jsonString), &host)
	assert.Nil(t, err)
	assert.Equal(t, types.Bool(true), host.Enabled)
	assert.Equal(t, "srv01", host.Hostname)
	assert.Equal(t, "dev.sys-int.de", host.Domain)
	assert.Equal(t, Rr("MX"), host.Rr)
	assert.Equal(t, "srv01.dev.sys-int.de", host.Mx)
	assert.Equal(t, MxPrio(10), host.Mxprio)
	assert.Equal(t, "srv01 - MX", host.Description)
	assert.Equal(t, "server01", host.Server)
}

func TestMxPrioToInt(t *testing.T) {
	cut := MxPrio(101)
	assert.Equal(t, 101, cut.Int())
}

func TestRrToString(t *testing.T) {
	cut := Rr("AAAA")
	assert.Equal(t, "AAAA", cut.String())
}
