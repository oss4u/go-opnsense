package unbound

import (
	"encoding/json"
	"fmt"
	"github.com/oss4u/go-opnsense/opnsense"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type MySuite struct{ suite.Suite }

func TestCreateUpdateDelete(t *testing.T) {
	if os.Getenv("OPNSENSE_ADDRESS") == "" {
		t.Skip("Missing credentials")
	}
	api := opnsense.GetOpnSenseClient("", "", "")
	overrides := Get_HostOverrides(api)
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

func TestToJson(t *testing.T) {
	//api := opnsense.GetOpnSenseClient("", "", "")
	//overrides := Get_HostOverrides(api)
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
	result, _ := json.Marshal(host)
	fmt.Print(string(result))

}

func TestFromJsonToOverridesHost(t *testing.T) {
	json_string := "{\"host\":{\"enabled\":\"1\",\"hostname\":\"srv01\",\"domain\":\"dev.sys-int.de\",\"rr\":{\"A\":{\"value\":\"A (IPv4 address)\",\"selected\":0},\"AAAA\":{\"value\":\"AAAA (IPv6 address)\",\"selected\":0},\"MX\":{\"value\":\"MX (Mail server)\",\"selected\":1}},\"mxprio\":\"10\",\"mx\":\"srv01.dev.sys-int.de\",\"server\":\"server01\",\"description\":\"srv01 - MX\"}}"
	host := OverridesHost{}
	err := json.Unmarshal([]byte(json_string), &host)
	assert.Nil(t, err)
	assert.NotNil(t, host.Host)
	assert.Equal(t, host.Host.Hostname, "srv01")
	assert.Equal(t, "srv01", host.Host.Hostname)
	assert.Equal(t, "dev.sys-int.de", host.Host.Domain)
	assert.Equal(t, "MX", host.Host.Rr)
	assert.Equal(t, "srv01.dev.sys-int.de", host.Host.Mx)
	assert.Equal(t, 10, host.Host.Mxprio)
	assert.Equal(t, "srv01 - MX", host.Host.Description)
	assert.Equal(t, "server01", host.Host.Server)
}

func TestFromJsonToOverridesHostDetails(t *testing.T) {
	json_string := "{\"enabled\":\"1\",\"hostname\":\"srv01\",\"domain\":\"dev.sys-int.de\",\"rr\":{\"A\":{\"value\":\"A (IPv4 address)\",\"selected\":0},\"AAAA\":{\"value\":\"AAAA (IPv6 address)\",\"selected\":0},\"MX\":{\"value\":\"MX (Mail server)\",\"selected\":1}},\"mxprio\":\"10\",\"mx\":\"srv01.dev.sys-int.de\",\"server\":\"server01\",\"description\":\"srv01 - MX\"}"
	host := OverridesHostDetails{}
	err := json.Unmarshal([]byte(json_string), &host)
	assert.Nil(t, err)
	assert.Equal(t, "srv01", host.Hostname)
	assert.Equal(t, "dev.sys-int.de", host.Domain)
	assert.Equal(t, "MX", host.Rr)
	assert.Equal(t, "srv01.dev.sys-int.de", host.Mx)
	assert.Equal(t, 10, host.Mxprio)
	assert.Equal(t, "srv01 - MX", host.Description)
	assert.Equal(t, "server01", host.Server)
}
