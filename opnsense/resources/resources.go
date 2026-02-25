package resources

import (
	"encoding/json"

	"github.com/oss4u/go-opnsense/opnsense"
)

type MutationResult struct {
	Result string `json:"result"`
	UUID   string `json:"uuid,omitempty"`
}

type API struct {
	api        *opnsense.OpnSenseApi
	module     string
	controller string
}

func New(api *opnsense.OpnSenseApi, module string, controller string) API {
	return API{
		api:        api,
		module:     module,
		controller: controller,
	}
}

func NewCore(api *opnsense.OpnSenseApi, controller string) API {
	return New(api, "core", controller)
}

func NewPlugin(api *opnsense.OpnSenseApi, plugin string, controller string) API {
	return New(api, plugin, controller)
}

func encodePayload(payload any) (string, error) {
	if payload == nil {
		return `{}`, nil
	}
	if raw, ok := payload.(string); ok {
		if raw == "" {
			return `{}`, nil
		}
		return raw, nil
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

func parseMutationResult(raw string) (MutationResult, error) {
	result := MutationResult{}
	err := json.Unmarshal([]byte(raw), &result)
	return result, err
}

func (r API) Add(payload any) (MutationResult, error) {
	data, err := encodePayload(payload)
	if err != nil {
		return MutationResult{}, err
	}
	raw, err := r.api.ModifyingRequest(r.module, r.controller, "add", data, []string{})
	if err != nil {
		return MutationResult{}, err
	}
	return parseMutationResult(raw)
}

func (r API) Get(uuid string) (string, int, error) {
	params := []string{}
	if uuid != "" {
		params = append(params, uuid)
	}
	return r.api.NonModifyingRequest(r.module, r.controller, "get", params)
}

func (r API) Search(payload any) (string, error) {
	data, err := encodePayload(payload)
	if err != nil {
		return "", err
	}
	return r.api.ModifyingRequest(r.module, r.controller, "search", data, []string{})
}

func (r API) Set(uuid string, payload any) (MutationResult, error) {
	data, err := encodePayload(payload)
	if err != nil {
		return MutationResult{}, err
	}
	params := []string{}
	if uuid != "" {
		params = append(params, uuid)
	}
	raw, err := r.api.ModifyingRequest(r.module, r.controller, "set", data, params)
	if err != nil {
		return MutationResult{}, err
	}
	return parseMutationResult(raw)
}

func (r API) Delete(uuid string) (MutationResult, error) {
	params := []string{}
	if uuid != "" {
		params = append(params, uuid)
	}
	raw, err := r.api.ModifyingRequest(r.module, r.controller, "del", `{}`, params)
	if err != nil {
		return MutationResult{}, err
	}
	return parseMutationResult(raw)
}

func (r API) Toggle(uuid string, enabled *bool) (MutationResult, error) {
	params := []string{}
	if uuid != "" {
		params = append(params, uuid)
	}
	if enabled != nil {
		if *enabled {
			params = append(params, "1")
		} else {
			params = append(params, "0")
		}
	}
	raw, err := r.api.ModifyingRequest(r.module, r.controller, "toggle", `{}`, params)
	if err != nil {
		return MutationResult{}, err
	}
	return parseMutationResult(raw)
}
