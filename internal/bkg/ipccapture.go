package bkg

import (
	"encoding/json"
	"fmt"
)

type CaptureParams struct {
	Text string `json:"text"`
}

func (s *Service) handleCapture(params json.RawMessage) Response {
	var p CaptureParams
	err := json.Unmarshal(params, &p)
	if err != nil {
		return Response{
			Success: false,
			Error:   fmt.Sprintf("invalid parameters: %v", err),
		}
	}

	node, err := s.svc.Capture(p.Text)
	if err != nil {
		return Response{
			Success: false,
			Error:   fmt.Sprintf("error capturing: %v", err),
		}
	}

	nodeJSON, err := json.Marshal(node)
	if err != nil {
		return Response{
			Success: false,
			Error:   fmt.Sprintf("error marshaling result: %v", err),
		}
	}

	return Response{
		Success: true,
		Data:    nodeJSON,
	}
}
