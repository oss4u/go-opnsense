package overrides

import (
	"encoding/json"
	"fmt"
	"github.com/oss4u/go-opnsense/opnsense"
)

type OverridesHostsApi struct {
	api        *opnsense.OpnSenseApi
	module     string
	controller string
}

type Result struct {
	Result string `json:"result"`
	Uuid   string `json:"uuid"`
}

func (o OverridesHostsApi) CreateOverridesHost(host *OverridesHost) (*OverridesHost, error) {
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

func (o OverridesHostsApi) ReadOverridesHost(uuid string) (*OverridesHost, error) {
	param := []string{}
	param = append(param, uuid)
	result, retCode, err := o.api.NonModifyingRequest(o.module, o.controller, "getHostOverride", param)
	if retCode == 200 {
		if result == `[]` {
			return nil, err
		}
		host := OverridesHost{}
		json.Unmarshal([]byte(result), &host)
		return &host, err
	} else {
		return nil, err
	}
}

func (o OverridesHostsApi) UpdateOverridesHost(host *OverridesHost) (*OverridesHost, error) {
	params := []string{}
	params = append(params, host.Host.GetUUID())
	data, err := json.Marshal(host)
	if err != nil {
		return nil, err
	}
	o.api.ModifyingRequest(o.module, o.controller, "setHostOverride", string(data), params)
	return host, nil
}

func (o OverridesHostsApi) DeleteOverridesHost(host *OverridesHost) error {
	return o.DeleteByIDOverridesHost(host.Host.GetUUID())
}

func (o OverridesHostsApi) DeleteByIDOverridesHost(uuid string) error {
	params := []string{}
	params = append(params, uuid)
	_, err := o.api.ModifyingRequest(o.module, o.controller, "delHostOverride", "", params)
	return err
}
