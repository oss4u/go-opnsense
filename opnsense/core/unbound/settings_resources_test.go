package unbound

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oss4u/go-opnsense/opnsense"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type recordedRequest struct {
	method string
	path   string
	body   string
}

func TestSettingsAPI_ResourceMethods_UseExpectedEndpoints(t *testing.T) {
	recorded := make([]recordedRequest, 0, 9)

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, readErr := io.ReadAll(r.Body)
		require.NoError(t, readErr)

		recorded = append(recorded, recordedRequest{
			method: r.Method,
			path:   r.URL.Path,
			body:   string(bodyBytes),
		})

		switch r.URL.Path {
		case "/api/unbound/settings/get_host_override/u-1":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"host":{"hostname":"srv01"}}`))
			return
		case "/api/unbound/settings/get_host_override":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"host":{}}`))
			return
		default:
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"result":"ok","uuid":"u"}`))
			return
		}
	}))
	defer server.Close()

	api := opnsense.NewOpnSenseClient(server.URL, "test-key", "test-secret")
	settings := GetSettingsAPI(api)

	addResult, err := settings.Add(SettingsResourceHostOverride, map[string]any{"host": map[string]any{"hostname": "srv01"}})
	require.NoError(t, err)
	assert.Equal(t, "ok", addResult.Result)

	resourceRaw, status, err := settings.GetResource(SettingsResourceHostOverride, "u-1")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Contains(t, resourceRaw, "srv01")

	resourceRaw, status, err = settings.GetResource(SettingsResourceHostOverride, "")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, `{"host":{}}`, resourceRaw)

	searchResult, err := settings.Search(SettingsResourceHostOverride, map[string]any{"current": 1, "rowCount": 20})
	require.NoError(t, err)
	assert.Equal(t, "ok", searchResult.Result)

	setResult, err := settings.SetResource(SettingsResourceHostOverride, "u-2", map[string]any{"host": map[string]any{"hostname": "srv02"}})
	require.NoError(t, err)
	assert.Equal(t, "ok", setResult.Result)

	deleteResult, err := settings.DeleteResource(SettingsResourceHostOverride, "u-3")
	require.NoError(t, err)
	assert.Equal(t, "ok", deleteResult.Result)

	toggleNilResult, err := settings.ToggleResource(SettingsResourceHostOverride, "u-4", nil)
	require.NoError(t, err)
	assert.Equal(t, "ok", toggleNilResult.Result)

	enabled := true
	toggleOnResult, err := settings.ToggleResource(SettingsResourceHostOverride, "u-5", &enabled)
	require.NoError(t, err)
	assert.Equal(t, "ok", toggleOnResult.Result)

	disabled := false
	toggleOffResult, err := settings.ToggleResource(SettingsResourceHostOverride, "u-6", &disabled)
	require.NoError(t, err)
	assert.Equal(t, "ok", toggleOffResult.Result)

	require.Len(t, recorded, 9)
	assert.Equal(t, recordedRequest{method: http.MethodPost, path: "/api/unbound/settings/add_host_override", body: `{"host":{"hostname":"srv01"}}`}, recorded[0])
	assert.Equal(t, recordedRequest{method: http.MethodGet, path: "/api/unbound/settings/get_host_override/u-1", body: ""}, recorded[1])
	assert.Equal(t, recordedRequest{method: http.MethodGet, path: "/api/unbound/settings/get_host_override", body: ""}, recorded[2])
	assert.Equal(t, recordedRequest{method: http.MethodPost, path: "/api/unbound/settings/search_host_override", body: `{"current":1,"rowCount":20}`}, recorded[3])
	assert.Equal(t, recordedRequest{method: http.MethodPost, path: "/api/unbound/settings/set_host_override/u-2", body: `{"host":{"hostname":"srv02"}}`}, recorded[4])
	assert.Equal(t, recordedRequest{method: http.MethodPost, path: "/api/unbound/settings/del_host_override/u-3", body: `{}`}, recorded[5])
	assert.Equal(t, recordedRequest{method: http.MethodPost, path: "/api/unbound/settings/toggle_host_override/u-4", body: `{}`}, recorded[6])
	assert.Equal(t, recordedRequest{method: http.MethodPost, path: "/api/unbound/settings/toggle_host_override/u-5/1", body: `{}`}, recorded[7])
	assert.Equal(t, recordedRequest{method: http.MethodPost, path: "/api/unbound/settings/toggle_host_override/u-6/0", body: `{}`}, recorded[8])
}

func TestSettingsAPI_ResourceMethods_ReturnErrorOnInvalidMutationResponse(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`not-json`))
	}))
	defer server.Close()

	api := opnsense.NewOpnSenseClient(server.URL, "test-key", "test-secret")
	settings := GetSettingsAPI(api)

	_, err := settings.Add(SettingsResourceHostOverride, map[string]any{"host": map[string]any{"hostname": "srv01"}})
	require.Error(t, err)
}
