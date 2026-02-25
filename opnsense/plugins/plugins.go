package plugins

import (
	"github.com/oss4u/go-opnsense/opnsense"
	"github.com/oss4u/go-opnsense/opnsense/plugins/caddy"
	pluginresources "github.com/oss4u/go-opnsense/opnsense/plugins/resources"
)

type API struct {
	api    *opnsense.OpnSenseApi
	plugin string
}

type Registry struct {
	api *opnsense.OpnSenseApi
}

func New(api *opnsense.OpnSenseApi, plugin string) API {
	return API{api: api, plugin: plugin}
}

func NewRegistry(api *opnsense.OpnSenseApi) Registry {
	return Registry{api: api}
}

func (r Registry) Plugin(plugin string) API {
	return New(r.api, plugin)
}

func (p API) Controller(controller string) pluginresources.API {
	return pluginresources.New(p.api, p.plugin, controller)
}

func (p API) Name() string {
	return p.plugin
}

func (p API) Caddy() caddy.API {
	return caddy.New(p.api)
}
