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

func (h OverridesHost) UnmarshalJSON(data []byte) error {
	type Alias OverridesHost
	aux := struct {
		Host OverridesHostDetails `json:"host"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

func (h OverridesHost) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Host OverridesHostDetails `json:"host"`
	}{
		Host: h.Host,
	})
}

type OverridesHostDetails struct {
	Uuid        string `json:"uuid,omitempty"`
	Enabled     bool   `json:"enabled"`
	Hostname    string `json:"hostname"`
	Domain      string `json:"domain"`
	Rr          string `json:"rr"`
	Description string `json:"description"`
	Mxprio      int    `json:"mxprio,omitempty"`
	Mx          string `json:"mx,omitempty"`
	Server      string `json:"server,omitempty"`
}

func (h OverridesHostDetails) ConvertToJson() string {
	b, _ := json.Marshal(h)
	return string(b)
}

func (h OverridesHostDetails) UnmarshalJSON(data []byte) error {
	aux := struct {
		Uuid        string `json:"uuid,omitempty"`
		Enabled     string `json:"enabled"`
		Hostname    string `json:"hostname"`
		Domain      string `json:"domain"`
		Rr          string `json:"rr"`
		Description string `json:"description"`
		Mxprio      int    `json:"mxprio,omitempty"`
		Mx          string `json:"mx,omitempty"`
		Server      string `json:"server,omitempty"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Enabled == "1" {
		h.Enabled = true
	} else {
		h.Enabled = false
	}
	return nil
}

func (h *OverridesHostDetails) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Uuid        string `json:"uuid,omitempty"`
		Enabled     string `json:"enabled"`
		Hostname    string `json:"hostname"`
		Domain      string `json:"domain"`
		Rr          string `json:"rr"`
		Description string `json:"description"`
		MxPrio      int    `json:"mxprio,omitempty"`
		Mx          string `json:"mx,omitempty"`
		Server      string `json:"server,omitempty"`
	}{
		Hostname:    h.Hostname,
		Domain:      h.Domain,
		Rr:          h.Rr,
		Description: h.Description,
		MxPrio:      h.Mxprio,
		Mx:          h.Mx,
		Server:      h.Server,
		Enabled:     IfThenElse(h.Enabled),
	})
}

func (h OverridesHostDetails) GetUUID() string {
	return h.Uuid
}

func (h OverridesHostDetails) SetUUID(uuid string) {
	h.Uuid = uuid
}

func IfThenElse(value bool) string {
	if value {
		return "1"
	}
	return "0"
}
