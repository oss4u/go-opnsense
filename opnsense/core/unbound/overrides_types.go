package unbound

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type ConvertableToJson interface {
	ConvertToJson() string
	GetUUID() string
	SetUUID(uuid string)
}

type OverridesHost struct {
	Host OverridesHostDetails `json:"host"`
}

func (h *OverridesHost) UnmarshalJSON(data []byte) error {
	fmt.Printf("OH %s\n", string(data))
	type Alias OverridesHost
	//aux := OverridesHost{}
	var x map[string]interface{}
	json.Unmarshal(data, &x)
	fmt.Println(x)
	aux := struct {
		Host OverridesHostDetails `json:"host"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	h.Host = aux.Host
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

func (h *OverridesHostDetails) UnmarshalJSON(data []byte) error {
	fmt.Printf("OHD %s\n", string(data))
	aux := struct {
		Uuid     string `json:"uuid,omitempty"`
		Enabled  string `json:"enabled"`
		Hostname string `json:"hostname"`
		Domain   string `json:"domain"`
		Rr       map[string]struct {
			Value    string `json:"value"`
			Selected int    `json:"selected"`
		} `json:"rr"`
		Description string `json:"description"`
		Mxprio      string `json:"mxprio,omitempty"`
		Mx          string `json:"mx,omitempty"`
		Server      string `json:"server,omitempty"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	h.Hostname = aux.Hostname
	h.Domain = aux.Domain
	h.Description = aux.Description
	h.Mxprio, _ = strconv.Atoi(aux.Mxprio)
	h.Mx = aux.Mx
	h.Server = aux.Server
	if aux.Enabled == "1" {
		h.Enabled = true
	} else {
		h.Enabled = false
	}
	for k, v := range aux.Rr {
		if v.Selected == 1 {
			h.Rr = k
		}
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
