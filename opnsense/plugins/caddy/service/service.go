package service

import (
	"encoding/json"

	"github.com/oss4u/go-opnsense/opnsense"
	pluginresources "github.com/oss4u/go-opnsense/opnsense/plugins/resources"
)

type SearchRequest struct {
	Current      int            `json:"current"`
	RowCount     int            `json:"rowCount"`
	Sort         map[string]any `json:"sort,omitempty"`
	SearchPhrase string         `json:"searchPhrase,omitempty"`
}

type SearchResponse struct {
	Total    int              `json:"total"`
	RowCount int              `json:"rowCount"`
	Current  int              `json:"current"`
	Rows     []map[string]any `json:"rows"`
}

type API struct {
	resource pluginresources.API
}

func New(api *opnsense.OpnSenseApi) API {
	return API{resource: pluginresources.New(api, "caddy", "service")}
}

func (s API) Resource() pluginresources.API {
	return s.resource
}

func (s API) Search(payload SearchRequest) (SearchResponse, error) {
	raw, err := s.resource.Search(payload)
	if err != nil {
		return SearchResponse{}, err
	}

	response := SearchResponse{}
	if err := json.Unmarshal([]byte(raw), &response); err != nil {
		return SearchResponse{}, err
	}

	return response, nil
}

func (s API) Add(payload any) (pluginresources.MutationResult, error) {
	return s.resource.Add(payload)
}

func (s API) Get(uuid string) (string, int, error) {
	return s.resource.Get(uuid)
}

func (s API) Set(uuid string, payload any) (pluginresources.MutationResult, error) {
	return s.resource.Set(uuid, payload)
}

func (s API) Delete(uuid string) (pluginresources.MutationResult, error) {
	return s.resource.Delete(uuid)
}

func (s API) Toggle(uuid string, enabled *bool) (pluginresources.MutationResult, error) {
	return s.resource.Toggle(uuid, enabled)
}
