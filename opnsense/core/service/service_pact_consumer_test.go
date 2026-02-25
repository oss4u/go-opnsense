//go:build pact
// +build pact

package service_test

import (
	"fmt"
	"testing"

	"github.com/oss4u/go-opnsense/opnsense"
	coreservice "github.com/oss4u/go-opnsense/opnsense/core/service"
	"github.com/pact-foundation/pact-go/v2/consumer"
	"github.com/pact-foundation/pact-go/v2/matchers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCoreServicePactConsumer_Search(t *testing.T) {
	logDir := t.TempDir()
	pactDir := t.TempDir()

	pact, err := consumer.NewV2Pact(consumer.MockHTTPProviderConfig{
		Consumer: "go-opnsense-core-service-consumer",
		Provider: "opnsense-core-service-provider",
		LogDir:   logDir,
		PactDir:  pactDir,
	})
	require.NoError(t, err)

	err = pact.AddInteraction().
		Given("core service resources can be searched").
		UponReceiving("a typed core/service search request").
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
			serviceAPI := coreservice.New(api)

			response, callErr := serviceAPI.Search(coreservice.SearchRequest{Current: 1, RowCount: 7, SearchPhrase: ""})
			require.NoError(t, callErr)
			require.Len(t, response.Rows, 1)
			assert.Equal(t, "configd", response.Rows[0].ID)
			return nil
		})

	require.NoError(t, err)
}
