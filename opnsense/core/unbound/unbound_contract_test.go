package unbound

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/oss4u/go-opnsense/opnsense"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type capturedRequest struct {
	Method string
	Path   string
}

func newContractServer() (*httptest.Server, *capturedRequest) {
	captured := &capturedRequest{}
	mu := &sync.Mutex{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		captured.Method = r.Method
		captured.Path = r.URL.Path
		mu.Unlock()
		_, _ = io.WriteString(w, `{"result":"ok","uuid":"u-1"}`)
	}))

	return server, captured
}

func TestUnboundContract_DiagnosticsOverviewService_FromDocs(t *testing.T) {
	server, captured := newContractServer()
	defer server.Close()

	api := opnsense.NewOpnSenseClient(server.URL, "test-key", "test-secret")
	client := New(api)

	tests := []struct {
		name           string
		expectedMethod string
		expectedPath   string
		call           func()
	}{
		{name: "diagnostics dumpcache", expectedMethod: http.MethodGet, expectedPath: "/api/unbound/diagnostics/dumpcache", call: func() { _, _, _ = client.Diagnostics.Dumpcache() }},
		{name: "diagnostics dumpinfra", expectedMethod: http.MethodGet, expectedPath: "/api/unbound/diagnostics/dumpinfra", call: func() { _, _, _ = client.Diagnostics.Dumpinfra() }},
		{name: "diagnostics listinsecure", expectedMethod: http.MethodGet, expectedPath: "/api/unbound/diagnostics/listinsecure", call: func() { _, _, _ = client.Diagnostics.Listinsecure() }},
		{name: "diagnostics listlocaldata", expectedMethod: http.MethodGet, expectedPath: "/api/unbound/diagnostics/listlocaldata", call: func() { _, _, _ = client.Diagnostics.Listlocaldata() }},
		{name: "diagnostics listlocalzones", expectedMethod: http.MethodGet, expectedPath: "/api/unbound/diagnostics/listlocalzones", call: func() { _, _, _ = client.Diagnostics.Listlocalzones() }},
		{name: "diagnostics stats", expectedMethod: http.MethodGet, expectedPath: "/api/unbound/diagnostics/stats", call: func() { _, _, _ = client.Diagnostics.Stats() }},
		{name: "diagnostics test_blocklist", expectedMethod: http.MethodPost, expectedPath: "/api/unbound/diagnostics/test_blocklist", call: func() { _, _ = client.Diagnostics.TestBlocklist(map[string]string{"force": "1"}) }},
		{name: "overview rolling", expectedMethod: http.MethodGet, expectedPath: "/api/unbound/overview/_rolling/1h/5", call: func() { _, _, _ = client.Overview.Rolling("1h", 5) }},
		{name: "overview get_policies", expectedMethod: http.MethodGet, expectedPath: "/api/unbound/overview/get_policies/u-1", call: func() { _, _, _ = client.Overview.GetPolicies("u-1") }},
		{name: "overview is_block_list_enabled", expectedMethod: http.MethodGet, expectedPath: "/api/unbound/overview/is_block_list_enabled", call: func() { _, _, _ = client.Overview.IsBlockListEnabled() }},
		{name: "overview is_enabled", expectedMethod: http.MethodGet, expectedPath: "/api/unbound/overview/is_enabled", call: func() { _, _, _ = client.Overview.IsEnabled() }},
		{name: "overview search_queries", expectedMethod: http.MethodGet, expectedPath: "/api/unbound/overview/search_queries", call: func() { _, _, _ = client.Overview.SearchQueries() }},
		{name: "overview totals", expectedMethod: http.MethodGet, expectedPath: "/api/unbound/overview/totals/100", call: func() { _, _, _ = client.Overview.Totals("100") }},
		{name: "service dnsbl", expectedMethod: http.MethodGet, expectedPath: "/api/unbound/service/dnsbl", call: func() { _, _, _ = client.Service.Dnsbl() }},
		{name: "service reconfigure", expectedMethod: http.MethodPost, expectedPath: "/api/unbound/service/reconfigure", call: func() { _, _ = client.Service.Reconfigure() }},
		{name: "service reconfigure_general", expectedMethod: http.MethodGet, expectedPath: "/api/unbound/service/reconfigure_general", call: func() { _, _, _ = client.Service.ReconfigureGeneral() }},
		{name: "service restart", expectedMethod: http.MethodPost, expectedPath: "/api/unbound/service/restart", call: func() { _, _ = client.Service.Restart() }},
		{name: "service start", expectedMethod: http.MethodPost, expectedPath: "/api/unbound/service/start", call: func() { _, _ = client.Service.Start() }},
		{name: "service status", expectedMethod: http.MethodGet, expectedPath: "/api/unbound/service/status", call: func() { _, _, _ = client.Service.Status() }},
		{name: "service stop", expectedMethod: http.MethodPost, expectedPath: "/api/unbound/service/stop", call: func() { _, _ = client.Service.Stop() }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.call()
			assert.Equal(t, tt.expectedMethod, captured.Method)
			assert.Equal(t, tt.expectedPath, captured.Path)
		})
	}
}

func TestUnboundContract_SettingsResources_FromDocs(t *testing.T) {
	server, captured := newContractServer()
	defer server.Close()

	api := opnsense.NewOpnSenseClient(server.URL, "test-key", "test-secret")
	settings := New(api).Settings

	resources := []SettingsResource{
		SettingsResourceACL,
		SettingsResourceDNSBL,
		SettingsResourceForward,
		SettingsResourceHostAlias,
		SettingsResourceHostOverride,
	}

	for _, resource := range resources {
		t.Run(string(resource)+" add/get/set/del/search/toggle", func(t *testing.T) {
			_, err := settings.Add(resource, map[string]any{"example": true})
			require.NoError(t, err)
			assert.Equal(t, http.MethodPost, captured.Method)
			assert.Equal(t, "/api/unbound/settings/add_"+string(resource), captured.Path)

			_, _, err = settings.GetResource(resource, "u-1")
			require.NoError(t, err)
			assert.Equal(t, http.MethodGet, captured.Method)
			assert.Equal(t, "/api/unbound/settings/get_"+string(resource)+"/u-1", captured.Path)

			_, err = settings.SetResource(resource, "u-1", map[string]any{"example": "changed"})
			require.NoError(t, err)
			assert.Equal(t, http.MethodPost, captured.Method)
			assert.Equal(t, "/api/unbound/settings/set_"+string(resource)+"/u-1", captured.Path)

			_, err = settings.DeleteResource(resource, "u-1")
			require.NoError(t, err)
			assert.Equal(t, http.MethodPost, captured.Method)
			assert.Equal(t, "/api/unbound/settings/del_"+string(resource)+"/u-1", captured.Path)

			_, err = settings.Search(resource, map[string]any{"current": 1, "rowCount": 20})
			require.NoError(t, err)
			assert.Equal(t, http.MethodPost, captured.Method)
			assert.Equal(t, "/api/unbound/settings/search_"+string(resource), captured.Path)

			enabled := true
			_, err = settings.ToggleResource(resource, "u-1", &enabled)
			require.NoError(t, err)
			assert.Equal(t, http.MethodPost, captured.Method)
			assert.Equal(t, "/api/unbound/settings/toggle_"+string(resource)+"/u-1/1", captured.Path)
		})
	}
}

func TestUnboundContract_SettingsGeneral_FromDocs(t *testing.T) {
	server, captured := newContractServer()
	defer server.Close()

	api := opnsense.NewOpnSenseClient(server.URL, "test-key", "test-secret")
	settings := New(api).Settings

	_, _, err := settings.Get()
	require.NoError(t, err)
	assert.Equal(t, http.MethodGet, captured.Method)
	assert.Equal(t, "/api/unbound/settings/get", captured.Path)

	_, err = settings.Set(map[string]any{"general": map[string]string{"enable": "1"}})
	require.NoError(t, err)
	assert.Equal(t, http.MethodPost, captured.Method)
	assert.Equal(t, "/api/unbound/settings/set", captured.Path)

	_, _, err = settings.GetNameservers()
	require.NoError(t, err)
	assert.Equal(t, http.MethodGet, captured.Method)
	assert.Equal(t, "/api/unbound/settings/get_nameservers", captured.Path)

	_, err = settings.UpdateBlocklist()
	require.NoError(t, err)
	assert.Equal(t, http.MethodPost, captured.Method)
	assert.Equal(t, "/api/unbound/settings/update_blocklist", captured.Path)
}
