//go:build pact
// +build pact

package unbound

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

func TestUnboundPactConsumer_ServiceStatus(t *testing.T) {
	logDir := t.TempDir()
	pactDir := t.TempDir()

	pact, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "go-opnsense-unbound-consumer",
		Provider: "opnsense-unbound-provider",
		LogDir:   logDir,
		PactDir:  pactDir,
	})
	require.NoError(t, err)

	err = pact.AddInteraction().
		Given("unbound service status is available").
		UponReceiving("a request for unbound service status").
		WithRequest(consumer.Method("GET"), "/api/unbound/service/status").
		WillRespondWith(200, func(builder *consumer.V2ResponseBuilder) {
			builder.Header("Content-Type", matchers.String("application/json"))
			builder.JSONBody(map[string]any{
				"status": matchers.Like("running"),
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			baseURL := fmt.Sprintf("http://%s:%d", config.Host, config.Port)
			api := opnsense.NewOpnSenseClient(baseURL, "test-key", "test-secret")
			client := New(api)

			raw, statusCode, callErr := client.Service.Status()
			require.NoError(t, callErr)
			assert.Equal(t, 200, statusCode)

			response := map[string]string{}
			unmarshalErr := json.Unmarshal([]byte(raw), &response)
			require.NoError(t, unmarshalErr)
			assert.Equal(t, "running", response["status"])

			return nil
		})

	require.NoError(t, err)
}

func TestUnboundPactConsumer_SettingsSet(t *testing.T) {
	logDir := t.TempDir()
	pactDir := t.TempDir()

	pact, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "go-opnsense-unbound-consumer",
		Provider: "opnsense-unbound-provider",
		LogDir:   logDir,
		PactDir:  pactDir,
	})
	require.NoError(t, err)

	err = pact.AddInteraction().
		Given("unbound settings can be changed").
		UponReceiving("a request to update unbound settings").
		WithRequest(consumer.Method("POST"), "/api/unbound/settings/set", func(builder *consumer.V2RequestBuilder) {
			builder.Header("Content-Type", matchers.String("application/json"))
			builder.JSONBody(map[string]any{
				"general": map[string]any{
					"enable": matchers.Like("1"),
				},
			})
		}).
		WillRespondWith(200, func(builder *consumer.V2ResponseBuilder) {
			builder.Header("Content-Type", matchers.String("application/json"))
			builder.JSONBody(map[string]any{
				"result": matchers.Like("ok"),
				"uuid":   matchers.Like("u-1"),
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			baseURL := fmt.Sprintf("http://%s:%d", config.Host, config.Port)
			api := opnsense.NewOpnSenseClient(baseURL, "test-key", "test-secret")
			client := New(api)

			result, callErr := client.Settings.Set(map[string]any{
				"general": map[string]string{"enable": "1"},
			})
			require.NoError(t, callErr)
			assert.Equal(t, "ok", result.Result)
			assert.Equal(t, "u-1", result.UUID)

			return nil
		})

	require.NoError(t, err)
}
