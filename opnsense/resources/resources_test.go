package resources_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oss4u/go-opnsense/opnsense"
	coreresources "github.com/oss4u/go-opnsense/opnsense/core/resources"
	pluginresources "github.com/oss4u/go-opnsense/opnsense/plugins/resources"
	baseresources "github.com/oss4u/go-opnsense/opnsense/resources"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type recordedRequest struct {
	method string
	path   string
	body   string
}

func TestCoreResourcesAPI_UsesExpectedEndpoints(t *testing.T) {
	recorded := make([]recordedRequest, 0, 6)

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, readErr := io.ReadAll(r.Body)
		require.NoError(t, readErr)

		recorded = append(recorded, recordedRequest{
			method: r.Method,
			path:   r.URL.Path,
			body:   string(bodyBytes),
		})

		switch r.URL.Path {
		case "/api/core/service/get/svc-1":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":"svc-1"}`))
			return
		case "/api/core/service/search":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"total":1,"rows":[{"id":"svc-1"}]}`))
			return
		default:
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"result":"ok","uuid":"svc-1"}`))
			return
		}
	}))
	defer server.Close()

	api := opnsense.NewOpnSenseClient(server.URL, "test-key", "test-secret")
	resourceAPI := coreresources.New(api, "service")

	addResult, err := resourceAPI.Add(map[string]any{"name": "configd"})
	require.NoError(t, err)
	assert.Equal(t, "ok", addResult.Result)
	assert.Equal(t, "svc-1", addResult.UUID)

	getRaw, getStatus, err := resourceAPI.Get("svc-1")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, getStatus)
	assert.Equal(t, `{"id":"svc-1"}`, getRaw)

	searchRaw, err := resourceAPI.Search(map[string]any{"current": 1, "rowCount": 7, "searchPhrase": ""})
	require.NoError(t, err)
	assert.Contains(t, searchRaw, `"total":1`)

	setResult, err := resourceAPI.Set("svc-1", map[string]any{"name": "configd-updated"})
	require.NoError(t, err)
	assert.Equal(t, "ok", setResult.Result)

	deleteResult, err := resourceAPI.Delete("svc-1")
	require.NoError(t, err)
	assert.Equal(t, "ok", deleteResult.Result)

	enabled := true
	toggleResult, err := resourceAPI.Toggle("svc-1", &enabled)
	require.NoError(t, err)
	assert.Equal(t, "ok", toggleResult.Result)

	require.Len(t, recorded, 6)
	assert.Equal(t, recordedRequest{method: http.MethodPost, path: "/api/core/service/add", body: `{"name":"configd"}`}, recorded[0])
	assert.Equal(t, recordedRequest{method: http.MethodGet, path: "/api/core/service/get/svc-1", body: ""}, recorded[1])
	assert.Equal(t, recordedRequest{method: http.MethodPost, path: "/api/core/service/search", body: `{"current":1,"rowCount":7,"searchPhrase":""}`}, recorded[2])
	assert.Equal(t, recordedRequest{method: http.MethodPost, path: "/api/core/service/set/svc-1", body: `{"name":"configd-updated"}`}, recorded[3])
	assert.Equal(t, recordedRequest{method: http.MethodPost, path: "/api/core/service/del/svc-1", body: `{}`}, recorded[4])
	assert.Equal(t, recordedRequest{method: http.MethodPost, path: "/api/core/service/toggle/svc-1/1", body: `{}`}, recorded[5])
}

func TestPluginResourcesAPI_UsesExpectedEndpoints(t *testing.T) {
	recorded := make([]recordedRequest, 0, 4)

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, readErr := io.ReadAll(r.Body)
		require.NoError(t, readErr)

		recorded = append(recorded, recordedRequest{
			method: r.Method,
			path:   r.URL.Path,
			body:   string(bodyBytes),
		})

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"result":"ok","uuid":"plugin-1"}`))
	}))
	defer server.Close()

	api := opnsense.NewOpnSenseClient(server.URL, "test-key", "test-secret")
	resourceAPI := pluginresources.New(api, "caddy", "service")

	_, err := resourceAPI.Add(map[string]any{"name": "site-a"})
	require.NoError(t, err)

	_, err = resourceAPI.Set("plugin-1", map[string]any{"name": "site-b"})
	require.NoError(t, err)

	disabled := false
	_, err = resourceAPI.Toggle("plugin-1", &disabled)
	require.NoError(t, err)

	_, err = resourceAPI.Delete("plugin-1")
	require.NoError(t, err)

	require.Len(t, recorded, 4)
	assert.Equal(t, recordedRequest{method: http.MethodPost, path: "/api/caddy/service/add", body: `{"name":"site-a"}`}, recorded[0])
	assert.Equal(t, recordedRequest{method: http.MethodPost, path: "/api/caddy/service/set/plugin-1", body: `{"name":"site-b"}`}, recorded[1])
	assert.Equal(t, recordedRequest{method: http.MethodPost, path: "/api/caddy/service/toggle/plugin-1/0", body: `{}`}, recorded[2])
	assert.Equal(t, recordedRequest{method: http.MethodPost, path: "/api/caddy/service/del/plugin-1", body: `{}`}, recorded[3])
}

func TestResourcesAPI_ReturnsErrorOnInvalidMutationPayload(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`not-json`))
	}))
	defer server.Close()

	api := opnsense.NewOpnSenseClient(server.URL, "test-key", "test-secret")
	resourceAPI := baseresources.NewCore(api, "service")

	_, err := resourceAPI.Add(map[string]any{"name": "svc"})
	require.Error(t, err)
}
