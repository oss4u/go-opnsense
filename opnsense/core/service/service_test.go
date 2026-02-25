package service_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oss4u/go-opnsense/opnsense"
	coreservice "github.com/oss4u/go-opnsense/opnsense/core/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceAPI_Search_DecodesTypedResponse(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/core/service/search", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"total":1,"rowCount":7,"current":1,"rows":[{"id":"configd","running":1,"name":"configd"}]}`))
	}))
	defer server.Close()

	api := opnsense.NewOpnSenseClient(server.URL, "test-key", "test-secret")
	serviceAPI := coreservice.New(api)

	response, err := serviceAPI.Search(coreservice.SearchRequest{Current: 1, RowCount: 7})
	require.NoError(t, err)
	require.Len(t, response.Rows, 1)
	assert.Equal(t, 1, response.Total)
	assert.Equal(t, "configd", response.Rows[0].ID)
	assert.Equal(t, 1, response.Rows[0].Running)
}

func TestServiceAPI_Search_ReturnsUnmarshalErrorOnInvalidJSON(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`not-json`))
	}))
	defer server.Close()

	api := opnsense.NewOpnSenseClient(server.URL, "test-key", "test-secret")
	serviceAPI := coreservice.New(api)

	_, err := serviceAPI.Search(coreservice.SearchRequest{Current: 1, RowCount: 7})
	require.Error(t, err)
}
