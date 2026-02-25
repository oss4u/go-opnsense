package caddy

import (
	"github.com/oss4u/go-opnsense/opnsense"
	caddyservice "github.com/oss4u/go-opnsense/opnsense/plugins/caddy/service"
	pluginresources "github.com/oss4u/go-opnsense/opnsense/plugins/resources"
)

type API struct {
	api *opnsense.OpnSenseApi
}

func New(api *opnsense.OpnSenseApi) API {
	return API{api: api}
}

func (c API) Controller(controller string) pluginresources.API {
	return pluginresources.New(c.api, "caddy", controller)
}

func (c API) Service() caddyservice.API {
	return caddyservice.New(c.api)
}
