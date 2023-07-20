package unbound

import (
	"encoding/json"
	"fmt"
	"github.com/oss4u/go-opnsense/opnsense"
)

func Get_HostOverrides(api *opnsense.OpnSenseApi) Overrides {
	return Overrides{
		api:        api,
		module:     "unbound",
		controller: "settings",
	}
}

type Overrides struct {
	api        *opnsense.OpnSenseApi
	module     string
	controller string
}

type Result struct {
	Result string `json:"result"`
	Uuid   string `json:"uuid"`
}

func (o Overrides) Create(host *OverridesHost) (*OverridesHost, error) {
	data, err := json.Marshal(host)
	if err != nil {
		return nil, err
	}
	request, err := o.api.ModifyingRequest(o.module, o.controller, "addHostOverride", string(data), []string{})
	result := Result{}
	json.Unmarshal([]byte(request), &result)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return nil, err
	}
	host.Host.Uuid = result.Uuid
	return host, nil
}

func (o Overrides) Read(uuid string) {
	param := []string{}
	param = append(param, uuid)
	o.api.NonModifyingRequest(o.module, o.controller, "getHostOverride", param)
}

func (o Overrides) Update(host *OverridesHost) (*OverridesHost, error) {
	params := []string{}
	params = append(params, host.Host.GetUUID())
	data, err := json.Marshal(host)
	if err != nil {
		return nil, err
	}
	o.api.ModifyingRequest(o.module, o.controller, "setHostOverride", string(data), params)
	return host, nil
}

func (o Overrides) Delete(host *OverridesHost) {
	o.DeleteByID(host.Host.GetUUID())
}

func (o Overrides) DeleteByID(uuid string) {
	params := []string{}
	params = append(params, uuid)
	request, err := o.api.ModifyingRequest(o.module, o.controller, "delHostOverride", "", params)
	fmt.Printf("Result: %s\n", request)
	if err != nil {
		return
	}
}
