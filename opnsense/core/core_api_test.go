package core_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oss4u/go-opnsense/opnsense"
	"github.com/oss4u/go-opnsense/opnsense/core"
	coreservice "github.com/oss4u/go-opnsense/opnsense/core/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCoreAPI_Controller_UsesCoreModuleEndpoints(t *testing.T) {
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
		_, _ = w.Write([]byte(`{"total":1,"rows":[{"id":"svc-1"}]}`))
	}))
	defer server.Close()

	api := opnsense.NewOpnSenseClient(server.URL, "test-key", "test-secret")
	coreAPI := core.New(api)

	raw, err := coreAPI.Controller("service").Search(map[string]any{"current": 1, "rowCount": 7})
	require.NoError(t, err)
	assert.Contains(t, raw, `"total":1`)
	assert.Equal(t, http.MethodPost, method)
	assert.Equal(t, "/api/core/service/search", path)
	assert.Equal(t, `{"current":1,"rowCount":7}`, body)
}

func TestCoreAPI_Unbound_DelegatesToUnboundModule(t *testing.T) {
	var method string
	var path string

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"running"}`))
	}))
	defer server.Close()

	api := opnsense.NewOpnSenseClient(server.URL, "test-key", "test-secret")
	coreAPI := core.New(api)

	raw, statusCode, err := coreAPI.Unbound().Service.Status()
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Contains(t, raw, "running")
	assert.Equal(t, http.MethodGet, method)
	assert.Equal(t, "/api/unbound/service/status", path)
}

func TestCoreAPI_Service_TypedEntryPointUsesCoreServiceController(t *testing.T) {
	var method string
	var path string

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		path = r.URL.Path
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"total":1,"rowCount":7,"current":1,"rows":[{"id":"configd"}]}`))
	}))
	defer server.Close()

	api := opnsense.NewOpnSenseClient(server.URL, "test-key", "test-secret")
	coreAPI := core.New(api)

	response, err := coreAPI.Service().Search(coreservice.SearchRequest{Current: 1, RowCount: 7})
	require.NoError(t, err)
	require.Len(t, response.Rows, 1)
	assert.Equal(t, "configd", response.Rows[0].ID)
	assert.Equal(t, http.MethodPost, method)
	assert.Equal(t, "/api/core/service/search", path)
}
