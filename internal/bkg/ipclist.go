package bkg

import (
	"encoding/json"
	"fmt"

	"github.com/adrianpk/tyn/internal/model"
)

// ListParams represents parameters for the list command
type ListParams struct {
	Type   string   `json:"type,omitempty"`
	Tags   []string `json:"tags,omitempty"`
	Places []string `json:"places,omitempty"`
	Status string   `json:"status,omitempty"`
}

// handleList processes a list command
func (s *Service) handleList(params json.RawMessage) Response {
	var p ListParams
	err := json.Unmarshal(params, &p)
	if err != nil {
		return Response{
			Success: false,
			Error:   fmt.Sprintf("invalid parameters: %v", err),
		}
	}

	var filter model.Filter

	if p.Type != "" {
		filter.Type = p.Type
	}

	if len(p.Tags) > 0 {
		filter.Tags = p.Tags
	}

	if len(p.Places) > 0 {
		filter.Places = p.Places
	}

	if p.Status != "" {
		filter.Status = p.Status
	}

	nodes, err := s.svc.List(filter)
	if err != nil {
		return Response{
			Success: false,
			Error:   fmt.Sprintf("error listing: %v", err),
		}
	}

	nodesJSON, err := json.Marshal(nodes)
	if err != nil {
		return Response{
			Success: false,
			Error:   fmt.Sprintf("error marshaling result: %v", err),
		}
	}

	return Response{
		Success: true,
		Data:    nodesJSON,
	}
}
