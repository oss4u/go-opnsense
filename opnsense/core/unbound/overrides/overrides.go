package overrides

import "github.com/oss4u/go-opnsense/opnsense"

func GetHostsOverrideApi(api *opnsense.OpnSenseApi) OverridesHostsApi {
	return OverridesHostsApi{
		api:        api,
		module:     "unbound",
		controller: "settings",
	}
}

func GetAliasesOverrideApi(api *opnsense.OpnSenseApi) OverridesAliasesApi {
	return OverridesAliasesApi{
		api:        api,
		module:     "unbound",
		controller: "settings",
	}
}
