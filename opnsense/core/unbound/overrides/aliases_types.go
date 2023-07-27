package overrides

import (
	"github.com/oss4u/go-opnsense/opnsense/types"
)

type OverridesAlias struct {
	Alias OverridesAliasDetails `json:"alias"`
}
type OverridesAliasDetails struct {
	Enabled     types.Bool `json:"enabled"`
	Host        string     `json:"host"`
	Hostname    string     `json:"hostname"`
	Domain      string     `json:"domain"`
	Description string     `json:"description"`
}
