package overrides

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/oss4u/go-opnsense/opnsense"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockApiState struct {
	hosts   map[string]string
	aliases map[string]OverridesAlias
}

func newMockOpnSenseServer(t *testing.T) *httptest.Server {
	state := &mockApiState{
		hosts:   map[string]string{},
		aliases: map[string]OverridesAlias{},
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user, pass, ok := r.BasicAuth(); !ok || user != "test-key" || pass != "test-secret" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"unauthorized"}`))
			return
		}

		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/unbound/settings/add_host_override":
			var host struct {
				Host map[string]any `json:"host"`
			}
			err := json.NewDecoder(r.Body).Decode(&host)
			require.NoError(t, err)
			hostname, _ := host.Host["hostname"].(string)
			state.hosts["host-1"] = hostname
			_, _ = w.Write([]byte(`{"result":"saved","uuid":"host-1"}`))
			return

		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/unbound/settings/get_host_override/"):
			uuid := strings.TrimPrefix(r.URL.Path, "/api/unbound/settings/get_host_override/")
			hostname, found := state.hosts[uuid]
			if !found {
				_, _ = w.Write([]byte(`[]`))
				return
			}
			payload := fmt.Sprintf(`{"host":{"enabled":"1","hostname":"%s","domain":"example.local","rr":{"A":{"value":"A","selected":1}},"description":"unit-test","server":"10.0.0.10"}}`, hostname)
			_, _ = w.Write([]byte(payload))
			return

		case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/api/unbound/settings/set_host_override/"):
			uuid := strings.TrimPrefix(r.URL.Path, "/api/unbound/settings/set_host_override/")
			var host struct {
				Host map[string]any `json:"host"`
			}
			err := json.NewDecoder(r.Body).Decode(&host)
			require.NoError(t, err)
			hostname, _ := host.Host["hostname"].(string)
			state.hosts[uuid] = hostname
			_, _ = w.Write([]byte(`{"result":"saved"}`))
			return

		case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/api/unbound/settings/del_host_override/"):
			uuid := strings.TrimPrefix(r.URL.Path, "/api/unbound/settings/del_host_override/")
			delete(state.hosts, uuid)
			_, _ = w.Write([]byte(`{"result":"deleted"}`))
			return

		case r.Method == http.MethodPost && r.URL.Path == "/api/unbound/settings/add_host_alias":
			var alias OverridesAlias
			err := json.NewDecoder(r.Body).Decode(&alias)
			require.NoError(t, err)
			state.aliases["alias-1"] = alias
			_, _ = w.Write([]byte(`{"result":"saved","uuid":"alias-1"}`))
			return

		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/api/unbound/settings/get_host_alias/"):
			uuid := strings.TrimPrefix(r.URL.Path, "/api/unbound/settings/get_host_alias/")
			alias, found := state.aliases[uuid]
			if !found {
				_, _ = w.Write([]byte(`[]`))
				return
			}
			payload, err := json.Marshal(alias)
			require.NoError(t, err)
			_, _ = w.Write(payload)
			return

		case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/api/unbound/settings/set_host_alias/"):
			uuid := strings.TrimPrefix(r.URL.Path, "/api/unbound/settings/set_host_alias/")
			var alias OverridesAlias
			err := json.NewDecoder(r.Body).Decode(&alias)
			require.NoError(t, err)
			state.aliases[uuid] = alias
			_, _ = w.Write([]byte(`{"result":"saved"}`))
			return

		case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/api/unbound/settings/del_host_alias/"):
			uuid := strings.TrimPrefix(r.URL.Path, "/api/unbound/settings/del_host_alias/")
			delete(state.aliases, uuid)
			_, _ = w.Write([]byte(`{"result":"deleted"}`))
			return
		}

		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, `{"error":"unsupported endpoint"}`)
	})

	return httptest.NewServer(handler)
}

func TestHostsOverrideApi_CRUD_WithMockedOpnSenseEndpoints(t *testing.T) {
	server := newMockOpnSenseServer(t)
	defer server.Close()

	api := opnsense.NewOpnSenseClient(server.URL, "test-key", "test-secret")
	hosts := GetHostsOverrideApi(api)

	host := &OverridesHost{Host: OverridesHostDetails{
		Enabled:     true,
		Hostname:    "unit-host",
		Domain:      "example.local",
		Rr:          "A",
		Description: "unit-test",
		Server:      "10.0.0.10",
	}}

	created, err := hosts.Create(host)
	require.NoError(t, err)
	assert.Equal(t, "host-1", created.Host.Uuid)

	created.Host.Hostname = "unit-host-updated"
	_, err = hosts.Update(created)
	require.NoError(t, err)

	readBack, err := hosts.Read("host-1")
	require.NoError(t, err)
	require.NotNil(t, readBack)
	assert.Equal(t, "unit-host-updated", readBack.Host.Hostname)

	err = hosts.DeleteByID("host-1")
	require.NoError(t, err)

	deleted, err := hosts.Read("host-1")
	require.NoError(t, err)
	assert.Nil(t, deleted)
}

func TestAliasesOverrideApi_CRUD_WithMockedOpnSenseEndpoints(t *testing.T) {
	server := newMockOpnSenseServer(t)
	defer server.Close()

	api := opnsense.NewOpnSenseClient(server.URL, "test-key", "test-secret")
	aliases := GetAliasesOverrideApi(api)

	alias := &OverridesAlias{Alias: OverridesAliasDetails{
		Enabled:     true,
		Host:        "host-1",
		Hostname:    "alias-host",
		Domain:      "example.local",
		Description: "unit-test-alias",
	}}

	created, err := aliases.Create(alias)
	require.NoError(t, err)
	assert.Equal(t, "alias-1", created.Alias.Uuid)

	created.Alias.Hostname = "alias-host-updated"
	_, err = aliases.Update(created)
	require.NoError(t, err)

	readBack, err := aliases.Read("alias-1")
	require.NoError(t, err)
	require.NotNil(t, readBack)
	assert.Equal(t, "alias-host-updated", readBack.Alias.Hostname)

	err = aliases.DeleteByID("alias-1")
	require.NoError(t, err)

	deleted, err := aliases.Read("alias-1")
	require.NoError(t, err)
	assert.Nil(t, deleted)
}
