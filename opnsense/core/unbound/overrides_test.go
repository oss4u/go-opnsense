package unbound

import (
	"encoding/json"
	"fmt"
	"github.com/oss4u/go-opnsense/opnsense"
	"testing"
)

func TestCreateUpdateDelete(t *testing.T) {
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
