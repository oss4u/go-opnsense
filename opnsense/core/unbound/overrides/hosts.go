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

func (o OverridesHostsApi) Create(host *OverridesHost) (*OverridesHost, error) {
	data, err := json.Marshal(host)
	if err != nil {
		return nil, err
	}
	request, err := o.api.ModifyingRequest(o.module, o.controller, "addHostOverride", string(data), []string{})
	if err != nil {
		return nil, fmt.Errorf("error creating host: %w", err)
	}
	var result Result
	if err := json.Unmarshal([]byte(request), &result); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	host.Host.Uuid = result.Uuid
	return host, nil
}

func (o OverridesHostsApi) Read(uuid string) (*OverridesHost, error) {
	params := []string{uuid}
	result, retCode, err := o.api.NonModifyingRequest(o.module, o.controller, "getHostOverride", params)
	if err != nil {
		return nil, fmt.Errorf("error reading host: %w", err)
	}
	if retCode != 200 || result == `[]` {
		return nil, fmt.Errorf("host not found or invalid response code: %d", retCode)
	}
	var host OverridesHost
	if err := json.Unmarshal([]byte(result), &host); err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}
	return &host, nil
}

func (o OverridesHostsApi) Update(host *OverridesHost) (*OverridesHost, error) {
	data, err := json.Marshal(host)
	if err != nil {
		return nil, err
	}
	params := []string{host.Host.GetUUID()}
	if _, err := o.api.ModifyingRequest(o.module, o.controller, "setHostOverride", string(data), params); err != nil {
		return nil, fmt.Errorf("error updating host: %w", err)
	}
	return host, nil
}

func (o OverridesHostsApi) Delete(host *OverridesHost) error {
	return o.DeleteByID(host.Host.GetUUID())
}

func (o OverridesHostsApi) DeleteByID(uuid string) error {
	params := []string{uuid}
	if _, err := o.api.ModifyingRequest(o.module, o.controller, "delHostOverride", "", params); err != nil {
		return fmt.Errorf("error deleting host: %w", err)
	}
	return nil
}
