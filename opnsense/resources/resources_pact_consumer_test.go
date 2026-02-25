//go:build pact
// +build pact

package resources

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/oss4u/go-opnsense/opnsense"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourcesPactConsumer_CoreSearch(t *testing.T) {
	logDir := t.TempDir()
	pactDir := t.TempDir()

	pact, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "go-opnsense-core-resources-consumer",
		Provider: "opnsense-core-resources-provider",
		LogDir:   logDir,
		PactDir:  pactDir,
	})
	require.NoError(t, err)

	err = pact.AddInteraction().
		Given("core services can be searched").
		UponReceiving("a request to search core service resources").
		WithRequest(consumer.Method("POST"), "/api/core/service/search", func(builder *consumer.V2RequestBuilder) {
			builder.Header("Content-Type", matchers.String("application/json"))
			builder.JSONBody(map[string]any{
				"current":      matchers.Like(1),
				"rowCount":     matchers.Like(7),
				"searchPhrase": matchers.Like(""),
			})
		}).
		WillRespondWith(200, func(builder *consumer.V2ResponseBuilder) {
			builder.Header("Content-Type", matchers.String("application/json"))
			builder.JSONBody(map[string]any{
				"total":    matchers.Like(1),
				"rowCount": matchers.Like(7),
				"current":  matchers.Like(1),
				"rows": matchers.EachLike(map[string]any{
					"id":      matchers.Like("configd"),
					"running": matchers.Like(1),
				}, 1),
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			baseURL := fmt.Sprintf("http://%s:%d", config.Host, config.Port)
			api := opnsense.NewOpnSenseClient(baseURL, "test-key", "test-secret")
			resourceAPI := NewCore(api, "service")

			raw, callErr := resourceAPI.Search(map[string]any{
				"current":      1,
				"rowCount":     7,
				"searchPhrase": "",
			})
			require.NoError(t, callErr)

			response := map[string]any{}
			unmarshalErr := json.Unmarshal([]byte(raw), &response)
			require.NoError(t, unmarshalErr)
			assert.Equal(t, float64(1), response["total"])

			return nil
		})

	require.NoError(t, err)
}

func TestResourcesPactConsumer_PluginToggle(t *testing.T) {
	logDir := t.TempDir()
	pactDir := t.TempDir()

	pact, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "go-opnsense-plugin-resources-consumer",
		Provider: "opnsense-plugin-resources-provider",
		LogDir:   logDir,
		PactDir:  pactDir,
	})
	require.NoError(t, err)

	err = pact.AddInteraction().
		Given("a plugin resource can be toggled").
		UponReceiving("a request to toggle plugin resource state").
		WithRequest(consumer.Method("POST"), "/api/caddy/service/toggle/plugin-1/1", func(builder *consumer.V2RequestBuilder) {
			builder.Header("Content-Type", matchers.String("application/json"))
			builder.JSONBody(map[string]any{})
		}).
		WillRespondWith(200, func(builder *consumer.V2ResponseBuilder) {
			builder.Header("Content-Type", matchers.String("application/json"))
			builder.JSONBody(map[string]any{
				"result": matchers.Like("ok"),
				"uuid":   matchers.Like("plugin-1"),
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			baseURL := fmt.Sprintf("http://%s:%d", config.Host, config.Port)
			api := opnsense.NewOpnSenseClient(baseURL, "test-key", "test-secret")
			resourceAPI := NewPlugin(api, "caddy", "service")

			enabled := true
			result, callErr := resourceAPI.Toggle("plugin-1", &enabled)
			require.NoError(t, callErr)
			assert.Equal(t, "ok", result.Result)
			assert.Equal(t, "plugin-1", result.UUID)
			return nil
		})

	require.NoError(t, err)
}
