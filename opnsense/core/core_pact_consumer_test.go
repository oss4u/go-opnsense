//go:build pact
// +build pact

package core_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/oss4u/go-opnsense/opnsense"
	"github.com/oss4u/go-opnsense/opnsense/core"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCorePactConsumer_ControllerSearch(t *testing.T) {
	logDir := t.TempDir()
	pactDir := t.TempDir()

	pact, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "go-opnsense-core-controller-consumer",
		Provider: "opnsense-core-controller-provider",
		LogDir:   logDir,
		PactDir:  pactDir,
	})
	require.NoError(t, err)

	err = pact.AddInteraction().
		Given("core service resources are searchable").
		UponReceiving("a request to search core/service resources via core API").
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
				"total": matchers.Like(1),
				"rows": matchers.EachLike(map[string]any{
					"id": matchers.Like("configd"),
				}, 1),
			})
		}).
		ExecuteTest(t, func(config consumer.MockServerConfig) error {
			baseURL := fmt.Sprintf("http://%s:%d", config.Host, config.Port)
			api := opnsense.NewOpnSenseClient(baseURL, "test-key", "test-secret")
			coreAPI := core.New(api)

			raw, callErr := coreAPI.Controller("service").Search(map[string]any{
				"current":      1,
				"rowCount":     7,
				"searchPhrase": "",
			})
			require.NoError(t, callErr)

			decoded := map[string]any{}
			unmarshalErr := json.Unmarshal([]byte(raw), &decoded)
			require.NoError(t, unmarshalErr)
			assert.Equal(t, float64(1), decoded["total"])
			return nil
		})

	require.NoError(t, err)
}
