package overrides

import (
	"encoding/json"
	"fmt"
	"github.com/oss4u/go-opnsense/opnsense"
)

type OverridesAliasesApi struct {
	api        *opnsense.OpnSenseApi
	module     string
	controller string
}

func (o OverridesAliasesApi) Create(host *OverridesHost) (*OverridesHost, error) {
	data, err := json.Marshal(host)
	if err != nil {
		return nil, err
	}
	request, err := o.api.ModifyingRequest(o.module, o.controller, "addHostAlias", string(data), []string{})
	result := Result{}
	json.Unmarshal([]byte(request), &result)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return nil, err
	}
	host.Host.Uuid = result.Uuid
	return host, nil
}

func (o OverridesAliasesApi) Read(uuid string) (*OverridesHost, error) {
	param := []string{}
	param = append(param, uuid)
	result, retCode, err := o.api.NonModifyingRequest(o.module, o.controller, "getHostAlias", param)
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

func (o OverridesAliasesApi) Update(host *OverridesHost) (*OverridesHost, error) {
	params := []string{}
	params = append(params, host.Host.GetUUID())
	data, err := json.Marshal(host)
	if err != nil {
		return nil, err
	}
	o.api.ModifyingRequest(o.module, o.controller, "setHostAlias", string(data), params)
	return host, nil
}

func (o OverridesAliasesApi) Delete(host *OverridesHost) error {
	return o.DeleteByID(host.Host.GetUUID())
}

func (o OverridesAliasesApi) DeleteByID(uuid string) error {
	params := []string{}
	params = append(params, uuid)
	_, err := o.api.ModifyingRequest(o.module, o.controller, "delHostAlias", "", params)
	return err
}
