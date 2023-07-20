package unbound

import "encoding/json"

type ConvertableToJson interface {
	ConvertToJson() string
	GetUUID() string
	SetUUID(uuid string)
}

type OverridesHost struct {
	Host OverridesHostDetails `json:"host"`
}

type OverridesHostDetails struct {
	Uuid        string `json:"uuid,omitempty"`
	Enabled     string `json:"enabled"`
	Hostname    string `json:"hostname"`
	Domain      string `json:"domain"`
	Rr          string `json:"rr"`
	Description string `json:"description"`
	Mxprio      string `json:"mxprio,omitempty"`
	Mx          int    `json:"mx,omitempty"`
	Server      string `json:"server,omitempty"`
}

func (h OverridesHostDetails) ConvertToJson() string {
	b, _ := json.Marshal(h)
	return string(b)
}

func (h OverridesHostDetails) GetUUID() string {
	return h.Uuid
}

func (h OverridesHostDetails) SetUUID(uuid string) {
	h.Uuid = uuid
}
