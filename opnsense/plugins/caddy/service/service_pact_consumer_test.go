//go:build pact
// +build pact

package service_test

import (
	"fmt"
	"testing"

	"github.com/oss4u/go-opnsense/opnsense"
	caddyservice "github.com/oss4u/go-opnsense/opnsense/plugins/caddy/service"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCaddyServicePactConsumer_Toggle(t *testing.T) {
	logDir := t.TempDir()
	pactDir := t.TempDir()

	pact, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "go-opnsense-caddy-service-consumer",
		Provider: "opnsense-caddy-service-provider",
		LogDir:   logDir,
		PactDir:  pactDir,
	})
	require.NoError(t, err)

	err = pact.AddInteraction().
		Given("caddy service resources can be toggled").
		UponReceiving("a typed caddy/service toggle request").
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
			serviceAPI := caddyservice.New(api)

			enabled := true
			result, callErr := serviceAPI.Toggle("plugin-1", &enabled)
			require.NoError(t, callErr)
			assert.Equal(t, "ok", result.Result)
			return nil
		})

	require.NoError(t, err)
}
