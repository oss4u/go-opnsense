//go:build manual
// +build manual

package opnsense_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/oss4u/go-opnsense/opnsense"
	"github.com/oss4u/go-opnsense/opnsense/core"
	coreservice "github.com/oss4u/go-opnsense/opnsense/core/service"
	"github.com/oss4u/go-opnsense/opnsense/core/unbound"
	"github.com/oss4u/go-opnsense/opnsense/core/unbound/overrides"
	"github.com/oss4u/go-opnsense/opnsense/plugins"
	caddyservice "github.com/oss4u/go-opnsense/opnsense/plugins/caddy/service"
	"github.com/stretchr/testify/require"
)

func manualClient(t *testing.T) *opnsense.OpnSenseApi {
	t.Helper()

	address := os.Getenv("OPNSENSE_ADDRESS")
	key := os.Getenv("OPNSENSE_KEY")
	secret := os.Getenv("OPNSENSE_SECRET")

	if address == "" || key == "" || secret == "" {
		t.Skip("manual tests require OPNSENSE_ADDRESS, OPNSENSE_KEY and OPNSENSE_SECRET")
	}

	return opnsense.NewOpnSenseClient(address, key, secret)
}

func requireJSON(t *testing.T, raw string) {
	t.Helper()

	if strings.HasPrefix(strings.TrimSpace(raw), "<") {
		t.Skip("manual test received HTML instead of JSON API response (check OPNSENSE_ADDRESS/API credentials)")
	}

	var decoded any
	require.NoError(t, json.Unmarshal([]byte(raw), &decoded))
}

func requireNoErrorOrSkipHTML(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		return
	}

	if strings.Contains(err.Error(), "invalid character '<'") {
		t.Skip("manual test received HTML instead of JSON API response (check OPNSENSE_ADDRESS/API credentials)")
	}

	require.NoError(t, err)
}

func TestManualResources_CoreAndTypedCoreServiceSearch(t *testing.T) {
	api := manualClient(t)

	coreAPI := core.New(api)

	raw, err := coreAPI.Controller("service").Search(map[string]any{
		"current":  1,
		"rowCount": 50,
	})
	requireNoErrorOrSkipHTML(t, err)
	requireJSON(t, raw)

	typedResponse, err := coreAPI.Service().Search(coreservice.SearchRequest{
		Current:  1,
		RowCount: 50,
	})
	requireNoErrorOrSkipHTML(t, err)
	require.GreaterOrEqual(t, typedResponse.Total, 0)
}

func TestManualResources_PluginControllerSearch(t *testing.T) {
	api := manualClient(t)

	pluginName := os.Getenv("OPNSENSE_MANUAL_PLUGIN")
	controllerName := os.Getenv("OPNSENSE_MANUAL_PLUGIN_CONTROLLER")
	if pluginName == "" || controllerName == "" {
		t.Skip("set OPNSENSE_MANUAL_PLUGIN and OPNSENSE_MANUAL_PLUGIN_CONTROLLER for plugin resource tests")
	}

	pluginAPI := plugins.New(api, pluginName)
	raw, err := pluginAPI.Controller(controllerName).Search(map[string]any{
		"current":  1,
		"rowCount": 50,
	})
	requireNoErrorOrSkipHTML(t, err)
	requireJSON(t, raw)

	if pluginName == "caddy" && controllerName == "service" {
		typedResponse, typedErr := pluginAPI.Caddy().Service().Search(caddyservice.SearchRequest{Current: 1, RowCount: 50})
		requireNoErrorOrSkipHTML(t, typedErr)
		require.GreaterOrEqual(t, typedResponse.Total, 0)
	}
}

func TestManualResources_UnboundSettingsResources(t *testing.T) {
	api := manualClient(t)

	settings := unbound.New(api).Settings
	resources := []unbound.SettingsResource{
		unbound.SettingsResourceACL,
		unbound.SettingsResourceDNSBL,
		unbound.SettingsResourceForward,
		unbound.SettingsResourceHostAlias,
		unbound.SettingsResourceHostOverride,
	}

	for _, resource := range resources {
		resource := resource
		t.Run(string(resource), func(t *testing.T) {
			searchResult, err := settings.Search(resource, map[string]any{"current": 1, "rowCount": 50})
			requireNoErrorOrSkipHTML(t, err)
			require.True(t, searchResult.Result == "" || searchResult.Result == "ok")

			resourceRaw, statusCode, getErr := settings.GetResource(resource, "")
			requireNoErrorOrSkipHTML(t, getErr)
			require.GreaterOrEqual(t, statusCode, 200)
			require.Less(t, statusCode, 300)
			requireJSON(t, resourceRaw)
		})
	}
}

func TestManualResources_UnboundOverrideReaders(t *testing.T) {
	api := manualClient(t)

	hostUUID := os.Getenv("OPNSENSE_MANUAL_HOST_OVERRIDE_UUID")
	if hostUUID != "" {
		host, err := overrides.GetHostsOverrideApi(api).Read(hostUUID)
		requireNoErrorOrSkipHTML(t, err)
		require.NotNil(t, host)
	}

	aliasUUID := os.Getenv("OPNSENSE_MANUAL_ALIAS_OVERRIDE_UUID")
	if aliasUUID != "" {
		alias, err := overrides.GetAliasesOverrideApi(api).Read(aliasUUID)
		requireNoErrorOrSkipHTML(t, err)
		require.NotNil(t, alias)
	}

	if hostUUID == "" && aliasUUID == "" {
		t.Skip("set OPNSENSE_MANUAL_HOST_OVERRIDE_UUID and/or OPNSENSE_MANUAL_ALIAS_OVERRIDE_UUID for override read tests")
	}
}
