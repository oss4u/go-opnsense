package overrides

import (
	"encoding/json"
	"fmt"
	"github.com/oss4u/go-opnsense/opnsense/types"
	"strconv"
	"strings"
)

type ConvertableToJson interface {
	ConvertToJson() string
	GetUUID() string
	SetUUID(uuid string)
}

type OverridesHost struct {
	Host OverridesHostDetails `json:"host"`
}

type Rr string

func (r *Rr) String() string {
	return *(*string)(r)
}

func (r *Rr) UnmarshalJSON(data []byte) error {
	aux := map[string]struct {
		Value    string `json:"value"`
		Selected int    `json:"selected"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	for k, v := range aux {
		if v.Selected == 1 {
			*r = Rr(k)
		}
	}
	return nil
}

func (r *Rr) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, *r)), nil
}

type MxPrio int

func (m *MxPrio) Int() int {
	return *(*int)(m)
}

func (m *MxPrio) UnmarshalJSON(data []byte) error {
	res := strings.Trim(string(data), "\"")
	value, _ := strconv.Atoi(res)
	*m = MxPrio(value)
	return nil
}

func (m MxPrio) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%d\"", m)), nil
}

type OverridesHostDetails struct {
	Uuid        string     `json:"-"`
	Enabled     types.Bool `json:"enabled"`
	Hostname    string     `json:"hostname"`
	Domain      string     `json:"domain"`
	Rr          Rr         `json:"rr"`
	Description string     `json:"description"`
	Mxprio      MxPrio     `json:"mxprio,omitempty"`
	Mx          string     `json:"mx,omitempty"`
	Server      string     `json:"server,omitempty"`
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
