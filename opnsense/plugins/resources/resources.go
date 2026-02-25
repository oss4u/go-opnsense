package resources

import (
	"github.com/oss4u/go-opnsense/opnsense"
	base "github.com/oss4u/go-opnsense/opnsense/resources"
)

type API = base.API
type MutationResult = base.MutationResult

func New(api *opnsense.OpnSenseApi, plugin string, controller string) API {
	return base.NewPlugin(api, plugin, controller)
}
