package overrides

import "github.com/oss4u/go-opnsense/opnsense"

type OverridesAliasesApi struct {
	api        *opnsense.OpnSenseApi
	module     string
	controller string
}
