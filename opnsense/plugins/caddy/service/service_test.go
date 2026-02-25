package service_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oss4u/go-opnsense/opnsense"
	caddyservice "github.com/oss4u/go-opnsense/opnsense/plugins/caddy/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCaddyServiceAPI_Search_DecodesResponse(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/caddy/service/search", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"total":1,"rowCount":10,"current":1,"rows":[{"id":"site-a","enabled":true}]}`))
	}))
	defer server.Close()

	api := opnsense.NewOpnSenseClient(server.URL, "test-key", "test-secret")
	serviceAPI := caddyservice.New(api)

	response, err := serviceAPI.Search(caddyservice.SearchRequest{Current: 1, RowCount: 10})
	require.NoError(t, err)
	require.Len(t, response.Rows, 1)
	assert.Equal(t, 1, response.Total)
	assert.Equal(t, "site-a", response.Rows[0]["id"])
}
