package unbound

import (
	"encoding/json"
	"fmt"

	"github.com/oss4u/go-opnsense/opnsense"
)

type SettingsResource string

const (
	SettingsResourceACL          SettingsResource = "acl"
	SettingsResourceDNSBL        SettingsResource = "dnsbl"
	SettingsResourceForward      SettingsResource = "forward"
	SettingsResourceHostAlias    SettingsResource = "host_alias"
	SettingsResourceHostOverride SettingsResource = "host_override"
)

type API struct {
	Diagnostics DiagnosticsAPI
	Overview    OverviewAPI
	Service     ServiceAPI
	Settings    SettingsAPI
}

type DiagnosticsAPI struct {
	api        *opnsense.OpnSenseApi
	module     string
	controller string
}

type OverviewAPI struct {
	api        *opnsense.OpnSenseApi
	module     string
	controller string
}

type ServiceAPI struct {
	api        *opnsense.OpnSenseApi
	module     string
	controller string
}

type SettingsAPI struct {
	api        *opnsense.OpnSenseApi
	module     string
	controller string
}

type MutationResult struct {
	Result string `json:"result"`
	UUID   string `json:"uuid,omitempty"`
}

func New(api *opnsense.OpnSenseApi) API {
	return API{
		Diagnostics: GetDiagnosticsAPI(api),
		Overview:    GetOverviewAPI(api),
		Service:     GetServiceAPI(api),
		Settings:    GetSettingsAPI(api),
	}
}

func GetDiagnosticsAPI(api *opnsense.OpnSenseApi) DiagnosticsAPI {
	return DiagnosticsAPI{api: api, module: "unbound", controller: "diagnostics"}
}

func GetOverviewAPI(api *opnsense.OpnSenseApi) OverviewAPI {
	return OverviewAPI{api: api, module: "unbound", controller: "overview"}
}

func GetServiceAPI(api *opnsense.OpnSenseApi) ServiceAPI {
	return ServiceAPI{api: api, module: "unbound", controller: "service"}
}

func GetSettingsAPI(api *opnsense.OpnSenseApi) SettingsAPI {
	return SettingsAPI{api: api, module: "unbound", controller: "settings"}
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
	b, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func parseMutationResult(raw string) (MutationResult, error) {
	result := MutationResult{}
	err := json.Unmarshal([]byte(raw), &result)
	return result, err
}

func (d DiagnosticsAPI) Dumpcache() (string, int, error) {
	return d.api.NonModifyingRequest(d.module, d.controller, "dumpcache", []string{})
}

func (d DiagnosticsAPI) Dumpinfra() (string, int, error) {
	return d.api.NonModifyingRequest(d.module, d.controller, "dumpinfra", []string{})
}

func (d DiagnosticsAPI) Listinsecure() (string, int, error) {
	return d.api.NonModifyingRequest(d.module, d.controller, "listinsecure", []string{})
}

func (d DiagnosticsAPI) Listlocaldata() (string, int, error) {
	return d.api.NonModifyingRequest(d.module, d.controller, "listlocaldata", []string{})
}

func (d DiagnosticsAPI) Listlocalzones() (string, int, error) {
	return d.api.NonModifyingRequest(d.module, d.controller, "listlocalzones", []string{})
}

func (d DiagnosticsAPI) Stats() (string, int, error) {
	return d.api.NonModifyingRequest(d.module, d.controller, "stats", []string{})
}

func (d DiagnosticsAPI) TestBlocklist(payload any) (MutationResult, error) {
	data, err := encodePayload(payload)
	if err != nil {
		return MutationResult{}, err
	}
	raw, err := d.api.ModifyingRequest(d.module, d.controller, "test_blocklist", data, []string{})
	if err != nil {
		return MutationResult{}, err
	}
	return parseMutationResult(raw)
}

func (o OverviewAPI) Rolling(timeperiod string, clients int) (string, int, error) {
	params := []string{}
	if timeperiod != "" {
		params = append(params, timeperiod)
		if clients > 0 {
			params = append(params, fmt.Sprintf("%d", clients))
		}
	}
	return o.api.NonModifyingRequest(o.module, o.controller, "_rolling", params)
}

func (o OverviewAPI) GetPolicies(uuid string) (string, int, error) {
	params := []string{}
	if uuid != "" {
		params = append(params, uuid)
	}
	return o.api.NonModifyingRequest(o.module, o.controller, "get_policies", params)
}

func (o OverviewAPI) IsBlockListEnabled() (string, int, error) {
	return o.api.NonModifyingRequest(o.module, o.controller, "is_block_list_enabled", []string{})
}

func (o OverviewAPI) IsEnabled() (string, int, error) {
	return o.api.NonModifyingRequest(o.module, o.controller, "is_enabled", []string{})
}

func (o OverviewAPI) SearchQueries() (string, int, error) {
	return o.api.NonModifyingRequest(o.module, o.controller, "search_queries", []string{})
}

func (o OverviewAPI) Totals(maximum string) (string, int, error) {
	params := []string{}
	if maximum != "" {
		params = append(params, maximum)
	}
	return o.api.NonModifyingRequest(o.module, o.controller, "totals", params)
}

func (s ServiceAPI) Dnsbl() (string, int, error) {
	return s.api.NonModifyingRequest(s.module, s.controller, "dnsbl", []string{})
}

func (s ServiceAPI) Reconfigure() (MutationResult, error) {
	raw, err := s.api.ModifyingRequest(s.module, s.controller, "reconfigure", `{}`, []string{})
	if err != nil {
		return MutationResult{}, err
	}
	return parseMutationResult(raw)
}

func (s ServiceAPI) ReconfigureGeneral() (string, int, error) {
	return s.api.NonModifyingRequest(s.module, s.controller, "reconfigure_general", []string{})
}

func (s ServiceAPI) Restart() (MutationResult, error) {
	raw, err := s.api.ModifyingRequest(s.module, s.controller, "restart", `{}`, []string{})
	if err != nil {
		return MutationResult{}, err
	}
	return parseMutationResult(raw)
}

func (s ServiceAPI) Start() (MutationResult, error) {
	raw, err := s.api.ModifyingRequest(s.module, s.controller, "start", `{}`, []string{})
	if err != nil {
		return MutationResult{}, err
	}
	return parseMutationResult(raw)
}

func (s ServiceAPI) Status() (string, int, error) {
	return s.api.NonModifyingRequest(s.module, s.controller, "status", []string{})
}

func (s ServiceAPI) Stop() (MutationResult, error) {
	raw, err := s.api.ModifyingRequest(s.module, s.controller, "stop", `{}`, []string{})
	if err != nil {
		return MutationResult{}, err
	}
	return parseMutationResult(raw)
}

func (s SettingsAPI) Get() (string, int, error) {
	return s.api.NonModifyingRequest(s.module, s.controller, "get", []string{})
}

func (s SettingsAPI) Set(payload any) (MutationResult, error) {
	data, err := encodePayload(payload)
	if err != nil {
		return MutationResult{}, err
	}
	raw, err := s.api.ModifyingRequest(s.module, s.controller, "set", data, []string{})
	if err != nil {
		return MutationResult{}, err
	}
	return parseMutationResult(raw)
}

func (s SettingsAPI) GetNameservers() (string, int, error) {
	return s.api.NonModifyingRequest(s.module, s.controller, "get_nameservers", []string{})
}

func (s SettingsAPI) UpdateBlocklist() (MutationResult, error) {
	raw, err := s.api.ModifyingRequest(s.module, s.controller, "update_blocklist", `{}`, []string{})
	if err != nil {
		return MutationResult{}, err
	}
	return parseMutationResult(raw)
}

func (s SettingsAPI) Add(resource SettingsResource, payload any) (MutationResult, error) {
	data, err := encodePayload(payload)
	if err != nil {
		return MutationResult{}, err
	}
	raw, err := s.api.ModifyingRequest(s.module, s.controller, "add_"+string(resource), data, []string{})
	if err != nil {
		return MutationResult{}, err
	}
	return parseMutationResult(raw)
}

func (s SettingsAPI) GetResource(resource SettingsResource, uuid string) (string, int, error) {
	params := []string{}
	if uuid != "" {
		params = append(params, uuid)
	}
	return s.api.NonModifyingRequest(s.module, s.controller, "get_"+string(resource), params)
}

func (s SettingsAPI) Search(resource SettingsResource, payload any) (MutationResult, error) {
	data, err := encodePayload(payload)
	if err != nil {
		return MutationResult{}, err
	}
	raw, err := s.api.ModifyingRequest(s.module, s.controller, "search_"+string(resource), data, []string{})
	if err != nil {
		return MutationResult{}, err
	}
	return parseMutationResult(raw)
}

func (s SettingsAPI) SetResource(resource SettingsResource, uuid string, payload any) (MutationResult, error) {
	data, err := encodePayload(payload)
	if err != nil {
		return MutationResult{}, err
	}
	raw, err := s.api.ModifyingRequest(s.module, s.controller, "set_"+string(resource), data, []string{uuid})
	if err != nil {
		return MutationResult{}, err
	}
	return parseMutationResult(raw)
}

func (s SettingsAPI) DeleteResource(resource SettingsResource, uuid string) (MutationResult, error) {
	raw, err := s.api.ModifyingRequest(s.module, s.controller, "del_"+string(resource), `{}`, []string{uuid})
	if err != nil {
		return MutationResult{}, err
	}
	return parseMutationResult(raw)
}

func (s SettingsAPI) ToggleResource(resource SettingsResource, uuid string, enabled *bool) (MutationResult, error) {
	params := []string{uuid}
	if enabled != nil {
		if *enabled {
			params = append(params, "1")
		} else {
			params = append(params, "0")
		}
	}
	raw, err := s.api.ModifyingRequest(s.module, s.controller, "toggle_"+string(resource), `{}`, params)
	if err != nil {
		return MutationResult{}, err
	}
	return parseMutationResult(raw)
}
