package overrides

import (
	"github.com/oss4u/go-opnsense/opnsense/types"
)

type OverridesAlias struct {
	Alias OverridesAliasDetails `json:"alias"`
}
type OverridesAliasDetails struct {
	Uuid        string     `json:"-"`
	Enabled     types.Bool `json:"enabled"`
	Host        string     `json:"host"`
	Hostname    string     `json:"hostname"`
	Domain      string     `json:"domain"`
	Description string     `json:"description"`
}

func NewOverridesAlias(enabled bool, host string, hostname string, domain string, description string) *OverridesAlias {
	return &OverridesAlias{
		Alias: OverridesAliasDetails{
			Enabled:     types.Bool(enabled),
			Host:        host,
			Hostname:    hostname,
			Domain:      domain,
			Description: description,
		},
	}
}
