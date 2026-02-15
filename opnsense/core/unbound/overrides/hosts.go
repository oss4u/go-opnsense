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
<<<<<<< Updated upstream
	request, err := o.api.ModifyingRequest(o.module, o.controller, "addHostOverride", string(data), []string{})
=======
	request, err := o.api.ModifyingRequest(o.module, o.controller, "add_host_override", string(data), []string{})
	result := Result{}
	json.Unmarshal([]byte(request), &result)
>>>>>>> Stashed changes
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
<<<<<<< Updated upstream
	params := []string{uuid}
	result, retCode, err := o.api.NonModifyingRequest(o.module, o.controller, "getHostOverride", params)
	if err != nil {
		return nil, fmt.Errorf("error reading host: %w", err)
=======
	param := []string{}
	param = append(param, uuid)
	result, retCode, err := o.api.NonModifyingRequest(o.module, o.controller, "get_host_override", param)
	if retCode == 200 {
		if result == `[]` {
			return nil, err
		}
		host := OverridesHost{}
		json.Unmarshal([]byte(result), &host)
		return &host, err
	} else {
		return nil, err
>>>>>>> Stashed changes
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
<<<<<<< Updated upstream
	params := []string{host.Host.GetUUID()}
	if _, err := o.api.ModifyingRequest(o.module, o.controller, "setHostOverride", string(data), params); err != nil {
		return nil, fmt.Errorf("error updating host: %w", err)
	}
=======
	o.api.ModifyingRequest(o.module, o.controller, "set_host_override", string(data), params)
>>>>>>> Stashed changes
	return host, nil
}

func (o OverridesHostsApi) Delete(host *OverridesHost) error {
	return o.DeleteByID(host.Host.GetUUID())
}

func (o OverridesHostsApi) DeleteByID(uuid string) error {
<<<<<<< Updated upstream
	params := []string{uuid}
	if _, err := o.api.ModifyingRequest(o.module, o.controller, "delHostOverride", "", params); err != nil {
		return fmt.Errorf("error deleting host: %w", err)
	}
	return nil
=======
	params := []string{}
	params = append(params, uuid)
	_, err := o.api.ModifyingRequest(o.module, o.controller, "del_host_override", "", params)
	return err
>>>>>>> Stashed changes
}
