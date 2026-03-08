package overrides

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/oss4u/go-opnsense/opnsense"
)

func extractUUIDFromResponse(raw string) string {
	if strings.TrimSpace(raw) == "" {
		return ""
	}

	var decoded any
	if err := json.Unmarshal([]byte(raw), &decoded); err != nil {
		return ""
	}

	return strings.TrimSpace(findUUIDInAny(decoded))
}

func findUUIDInAny(v any) string {
	switch typed := v.(type) {
	case map[string]any:
		if value, ok := typed["uuid"].(string); ok && strings.TrimSpace(value) != "" {
			return value
		}

		if result := typed["result"]; result != nil {
			if uuid := findUUIDInAny(result); uuid != "" {
				return uuid
			}
		}

		for _, child := range typed {
			if uuid := findUUIDInAny(child); uuid != "" {
				return uuid
			}
		}
	case []any:
		for _, child := range typed {
			if uuid := findUUIDInAny(child); uuid != "" {
				return uuid
			}
		}
	}

	return ""
}

func findHostOverrideUUIDByHostDomain(api *opnsense.OpnSenseApi, hostname, domain string) (string, error) {
	searchPayload := map[string]any{
		"current":  1,
		"rowCount": 500,
	}

	payloadRaw, err := json.Marshal(searchPayload)
	if err != nil {
		return "", err
	}

	raw, err := api.ModifyingRequest("unbound", "settings", "search_host_override", string(payloadRaw), []string{})
	if err != nil {
		return "", err
	}

	var decoded map[string]any
	if err := json.Unmarshal([]byte(raw), &decoded); err != nil {
		return "", err
	}

	rows, ok := decoded["rows"].([]any)
	if !ok {
		return "", nil
	}

	for _, row := range rows {
		rowMap, ok := row.(map[string]any)
		if !ok {
			continue
		}

		candidateHost, _ := rowMap["hostname"].(string)
		candidateDomain, _ := rowMap["domain"].(string)
		uuid, _ := rowMap["uuid"].(string)

		if strings.EqualFold(strings.TrimSpace(candidateHost), strings.TrimSpace(hostname)) &&
			strings.EqualFold(strings.TrimSpace(candidateDomain), strings.TrimSpace(domain)) {
			return strings.TrimSpace(uuid), nil
		}
	}

	return "", nil
}

func findHostAliasUUIDByHostDomain(api *opnsense.OpnSenseApi, hostname, domain string) (string, error) {
	searchPayload := map[string]any{
		"current":  1,
		"rowCount": 500,
	}

	payloadRaw, err := json.Marshal(searchPayload)
	if err != nil {
		return "", err
	}

	raw, err := api.ModifyingRequest("unbound", "settings", "search_host_alias", string(payloadRaw), []string{})
	if err != nil {
		return "", err
	}

	var decoded map[string]any
	if err := json.Unmarshal([]byte(raw), &decoded); err != nil {
		return "", err
	}

	rows, ok := decoded["rows"].([]any)
	if !ok {
		return "", nil
	}

	for _, row := range rows {
		rowMap, ok := row.(map[string]any)
		if !ok {
			continue
		}

		candidateHost, _ := rowMap["hostname"].(string)
		candidateDomain, _ := rowMap["domain"].(string)
		uuid, _ := rowMap["uuid"].(string)

		if strings.EqualFold(strings.TrimSpace(candidateHost), strings.TrimSpace(hostname)) &&
			strings.EqualFold(strings.TrimSpace(candidateDomain), strings.TrimSpace(domain)) {
			return strings.TrimSpace(uuid), nil
		}
	}

	return "", nil
}

func requireUUID(entity, rawResponse, parsedUUID string) error {
	if strings.TrimSpace(parsedUUID) != "" {
		return nil
	}

	return fmt.Errorf("%s create response did not contain a uuid: %s", entity, strings.TrimSpace(rawResponse))
}
