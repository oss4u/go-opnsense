package core

import (
	"github.com/oss4u/go-opnsense/opnsense"
	coreresources "github.com/oss4u/go-opnsense/opnsense/core/resources"
	coreservice "github.com/oss4u/go-opnsense/opnsense/core/service"
	"github.com/oss4u/go-opnsense/opnsense/core/unbound"
)

type API struct {
	api *opnsense.OpnSenseApi
}

func New(api *opnsense.OpnSenseApi) API {
	return API{api: api}
}

func (c API) Controller(controller string) coreresources.API {
	return coreresources.New(c.api, controller)
}

func (c API) Unbound() unbound.API {
	return unbound.New(c.api)
}

func (c API) Service() coreservice.API {
	return coreservice.New(c.api)
}
