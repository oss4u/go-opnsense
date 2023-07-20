package ifaces

import "github.com/oss4u/go-opnsense/opnsense/core/unbound"

type IOverrides interface {
	Create(host *unbound.OverridesHost) (*unbound.OverridesHost, error)
	Read(uuid string)
	Update(host *unbound.OverridesHost) (*unbound.OverridesHost, error)
	Delete(host unbound.OverridesHost)
	DeleteByID(uuid string)
}
