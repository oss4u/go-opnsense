package plugins_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oss4u/go-opnsense/opnsense"
	"github.com/oss4u/go-opnsense/opnsense/plugins"
	caddyservice "github.com/oss4u/go-opnsense/opnsense/plugins/caddy/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPluginsRegistry_ResolvesPluginAndController(t *testing.T) {
	var method string
	var path string
	var body string

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		method = r.Method
		path = r.URL.Path
		body = string(b)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result":"ok","uuid":"plugin-1"}`))
	}))
	defer server.Close()

	api := opnsense.NewOpnSenseClient(server.URL, "test-key", "test-secret")
	registry := plugins.NewRegistry(api)
	plugin := registry.Plugin("caddy")

	result, err := plugin.Controller("service").Add(map[string]any{"name": "site-a"})
	require.NoError(t, err)
	assert.Equal(t, "ok", result.Result)
	assert.Equal(t, "plugin-1", result.UUID)
	assert.Equal(t, "caddy", plugin.Name())
	assert.Equal(t, http.MethodPost, method)
	assert.Equal(t, "/api/caddy/service/add", path)
	assert.Equal(t, `{"name":"site-a"}`, body)
}

func TestPluginsAPI_Toggle_UsesPluginControllerEndpoint(t *testing.T) {
	var method string
	var path string

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result":"ok","uuid":"plugin-1"}`))
	}))
	defer server.Close()

	api := opnsense.NewOpnSenseClient(server.URL, "test-key", "test-secret")
	plugin := plugins.New(api, "caddy")
	enabled := true

	result, err := plugin.Controller("service").Toggle("plugin-1", &enabled)
	require.NoError(t, err)
	assert.Equal(t, "ok", result.Result)
	assert.Equal(t, "plugin-1", result.UUID)
	assert.Equal(t, http.MethodPost, method)
	assert.Equal(t, "/api/caddy/service/toggle/plugin-1/1", path)
}

func TestPluginsAPI_CaddyTypedEntryPoint_UsesCaddyServiceController(t *testing.T) {
	var method string
	var path string

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"total":1,"rowCount":10,"current":1,"rows":[{"id":"site-a"}]}`))
	}))
	defer server.Close()

	api := opnsense.NewOpnSenseClient(server.URL, "test-key", "test-secret")
	plugin := plugins.New(api, "caddy")

	response, err := plugin.Caddy().Service().Search(caddyservice.SearchRequest{Current: 1, RowCount: 10})
	require.NoError(t, err)
	require.Len(t, response.Rows, 1)
	assert.Equal(t, "site-a", response.Rows[0]["id"])
	assert.Equal(t, http.MethodPost, method)
	assert.Equal(t, "/api/caddy/service/search", path)
}
