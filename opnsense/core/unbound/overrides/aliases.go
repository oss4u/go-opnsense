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
	request, err := o.api.ModifyingRequest(o.module, o.controller, "add_host_alias", string(data), []string{})
	if err != nil {
		return nil, fmt.Errorf("error creating alias: %w", err)
	}
	var result Result
	if err := json.Unmarshal([]byte(request), &result); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	alias.Alias.Uuid = result.Uuid
	return alias, nil
}

func (o OverridesAliasesApi) Read(uuid string) (*OverridesAlias, error) {
	params := []string{uuid}
	result, retCode, err := o.api.NonModifyingRequest(o.module, o.controller, "get_host_alias", params)
	if err != nil {
		return nil, fmt.Errorf("error reading alias: %w", err)
	}
	if retCode == 200 && result == `[]` {
		return nil, nil
	}
	if retCode == 404 {
		return nil, nil
	}
	if retCode != 200 {
		return nil, fmt.Errorf("alias not found or invalid response code: %d", retCode)
	}
	var alias OverridesAlias
	if err := json.Unmarshal([]byte(result), &alias); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	return &alias, nil
}

func (o OverridesAliasesApi) Update(alias *OverridesAlias) (*OverridesAlias, error) {
	data, err := json.Marshal(alias)
	if err != nil {
		return nil, err
	}
	params := []string{alias.Alias.Uuid}
	if _, err := o.api.ModifyingRequest(o.module, o.controller, "set_host_alias", string(data), params); err != nil {
		return nil, fmt.Errorf("error updating alias: %w", err)
	}
	return alias, nil
}

func (o OverridesAliasesApi) Delete(alias *OverridesAlias) error {
	return o.DeleteByID(alias.Alias.Uuid)
}

func (o OverridesAliasesApi) DeleteByID(uuid string) error {
	params := []string{uuid}
	if _, err := o.api.ModifyingRequest(o.module, o.controller, "del_host_alias", "", params); err != nil {
		return fmt.Errorf("error deleting alias: %w", err)
	}
	return nil
}
