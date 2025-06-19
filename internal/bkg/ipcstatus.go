package bkg

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/adrianpk/tyn/internal/model"
)

type StatusParams struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Operation string `json:"operation"` // "set", "next", or "prev"
}

func (s *Service) handleStatus(p json.RawMessage) Response {
	var params StatusParams
	err := json.Unmarshal(p, &params)
	if err != nil {
		log.Printf("Error parsing status params: %v", err)
		return Response{Success: false, Error: fmt.Sprintf("error parsing params: %v", err)}
	}

	log.Printf("Status change requested: ID=%s, Status=%s, Operation=%s", params.ID, params.Status, params.Operation)

	ctx := context.Background()
	nodes, err := s.svc.Repo.List(ctx)
	if err != nil {
		log.Printf("Error listing nodes: %v", err)
		return Response{Success: false, Error: fmt.Sprintf("error listing nodes: %v", err)}
	}

	var targetNode model.Node
	var found bool

	for _, node := range nodes {
		if len(params.ID) <= 4 && node.ShortID() == params.ID {
			targetNode = node
			found = true
			log.Printf("Found node by ShortID: %s", node.ID)
			break
		} else if len(params.ID) > 4 && len(node.ID) >= len(params.ID) && node.ID[:len(params.ID)] == params.ID {
			targetNode = node
			found = true
			log.Printf("Found node by prefix: %s", node.ID)
			break
		}
	}

	if !found {
		log.Printf("Task with ID '%s' not found", params.ID)
		return Response{Success: false, Error: fmt.Sprintf("task with ID '%s' not found", params.ID)}
	}

	if targetNode.Type != "task" {
		log.Printf("Node with ID '%s' is not a task", params.ID)
		return Response{Success: false, Error: fmt.Sprintf("node with ID '%s' is not a task", params.ID)}
	}

	originalStatus := targetNode.Status
	var newStatus string

	switch params.Operation {
	case "set":
		if !model.ValidStatus(params.Status) {
			log.Printf("Invalid status: %s", params.Status)
			return Response{Success: false, Error: fmt.Sprintf("invalid status: %s", params.Status)}
		}
		newStatus = params.Status
	case "next":
		newStatus = model.NextStatus(targetNode.Status)
	case "prev":
		newStatus = model.PreviousStatus(targetNode.Status)
	default:
		log.Printf("Invalid operation: %s", params.Operation)
		return Response{Success: false, Error: fmt.Sprintf("invalid operation: %s", params.Operation)}
	}

	log.Printf("Changing status from '%s' to '%s'", originalStatus, newStatus)
	targetNode.Status = newStatus

	err = s.svc.Repo.Update(ctx, targetNode)
	if err != nil {
		log.Printf("Error updating task: %v", err)
		return Response{Success: false, Error: fmt.Sprintf("error updating task: %v", err)}
	}

	log.Printf("Status updated successfully: '%s' → '%s'", originalStatus, newStatus)

	result := struct {
		OriginalStatus string `json:"original_status"`
		NewStatus      string `json:"new_status"`
	}{
		OriginalStatus: originalStatus,
		NewStatus:      newStatus,
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		log.Printf("Error encoding result: %v", err)
		return Response{Success: false, Error: fmt.Sprintf("error encoding result: %v", err)}
	}

	return Response{Success: true, Data: resultJSON}
}
