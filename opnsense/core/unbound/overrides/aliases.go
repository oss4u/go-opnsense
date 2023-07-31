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

func (o OverridesAliasesApi) Create(alias *OverridesAlias) (*OverridesAlias, error) {
	data, err := json.Marshal(alias)
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
	alias.Alias.Uuid = result.Uuid
	return alias, nil
}

func (o OverridesAliasesApi) Read(uuid string) (*OverridesAlias, error) {
	param := []string{}
	param = append(param, uuid)
	result, retCode, err := o.api.NonModifyingRequest(o.module, o.controller, "getHostAlias", param)
	if retCode == 200 {
		if result == `[]` {
			return nil, err
		}
		host := OverridesAlias{}
		json.Unmarshal([]byte(result), &host)
		return &host, err
	} else {
		return nil, err
	}
}

func (o OverridesAliasesApi) Update(alias *OverridesAlias) (*OverridesAlias, error) {
	params := []string{}
	params = append(params, alias.Alias.Uuid)
	data, err := json.Marshal(alias)
	if err != nil {
		return nil, err
	}
	o.api.ModifyingRequest(o.module, o.controller, "setHostAlias", string(data), params)
	return alias, nil
}

func (o OverridesAliasesApi) Delete(alias *OverridesAlias) error {
	return o.DeleteByID(alias.Alias.Uuid)
}

func (o OverridesAliasesApi) DeleteByID(uuid string) error {
	params := []string{}
	params = append(params, uuid)
	_, err := o.api.ModifyingRequest(o.module, o.controller, "delHostAlias", "", params)
	return err
}
