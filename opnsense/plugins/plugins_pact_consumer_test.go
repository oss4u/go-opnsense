//go:build pact
// +build pact

package plugins_test

import (
	"fmt"
	"testing"

	"github.com/oss4u/go-opnsense/opnsense"
	"github.com/oss4u/go-opnsense/opnsense/plugins"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPluginsPactConsumer_ControllerToggle(t *testing.T) {
	logDir := t.TempDir()
	pactDir := t.TempDir()

	pact, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "go-opnsense-plugins-controller-consumer",
		Provider: "opnsense-plugins-controller-provider",
		LogDir:   logDir,
		PactDir:  pactDir,
	})
	require.NoError(t, err)

	err = pact.AddInteraction().
		Given("plugin resources can be toggled").
		UponReceiving("a request to toggle caddy/service resource via plugins API").
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
			plugin := plugins.New(api, "caddy")

			enabled := true
			result, callErr := plugin.Controller("service").Toggle("plugin-1", &enabled)
			require.NoError(t, callErr)
			assert.Equal(t, "ok", result.Result)
			assert.Equal(t, "plugin-1", result.UUID)
			return nil
		})

	require.NoError(t, err)
}
