package ifaces

import (
	"github.com/oss4u/go-opnsense/opnsense/core/unbound/overrides"
)

type IOverrides interface {
	Create(host *overrides.OverridesHost) (*overrides.OverridesHost, error)
	Read(uuid string)
	Update(host *overrides.OverridesHost) (*overrides.OverridesHost, error)
	Delete(host overrides.OverridesHost)
	DeleteByID(uuid string)
}
