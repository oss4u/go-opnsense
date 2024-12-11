package overrides

import "github.com/oss4u/go-opnsense/opnsense"

// GetHostsOverrideApi returns an instance of OverridesHostsApi with the provided OpnSenseApi.
func GetHostsOverrideApi(api *opnsense.OpnSenseApi) OverridesHostsApi {
	return OverridesHostsApi{
		api:        api,
		module:     "unbound",
		controller: "settings",
	}
}

// GetAliasesOverrideApi returns an instance of OverridesAliasesApi with the provided OpnSenseApi.
func GetAliasesOverrideApi(api *opnsense.OpnSenseApi) OverridesAliasesApi {
	return OverridesAliasesApi{
		api:        api,
		module:     "unbound",
		controller: "settings",
	}
}
