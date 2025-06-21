package bkg

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

func (s *Service) handlePlace(params json.RawMessage) Response {
	var placeParams PlaceCmdParams
	err := json.Unmarshal(params, &placeParams)
	if err != nil {
		return Response{
			Success: false,
			Error:   fmt.Sprintf("error unmarshaling place params: %v", err),
		}
	}

	log.Printf("Handling place operation: %s for task %s with places %v", placeParams.Operation, placeParams.ID, placeParams.Places)

	ctx := context.Background()
	task, err := s.svc.Repo.GetTaskByID(ctx, placeParams.ID)
	if err != nil {
		return Response{
			Success: false,
			Error:   fmt.Sprintf("error fetching task with ID %s: %v", placeParams.ID, err),
		}
	}

	var message string

	switch placeParams.Operation {
	case "add":
		for _, place := range placeParams.Places {
			exists := false
			for _, existingPlace := range task.Places {
				if existingPlace == place {
					exists = true
					break
				}
			}
			if !exists {
				task.Places = append(task.Places, place)
			}
		}
		message = fmt.Sprintf("Added places %v to task %s", placeParams.Places, placeParams.ID)

	case "remove":
		for _, placeToRemove := range placeParams.Places {
			newPlaces := []string{}
			for _, existingPlace := range task.Places {
				if existingPlace != placeToRemove {
					newPlaces = append(newPlaces, existingPlace)
				}
			}
			task.Places = newPlaces
		}
		message = fmt.Sprintf("Removed places %v from task %s", placeParams.Places, placeParams.ID)

	case "clear":
		task.Places = []string{}
		message = fmt.Sprintf("Cleared all places from task %s", placeParams.ID)

	default:
		return Response{
			Success: false,
			Error:   fmt.Sprintf("unknown place operation: %s", placeParams.Operation),
		}
	}

	err = s.svc.Repo.Update(ctx, task)
	if err != nil {
		return Response{
			Success: false,
			Error:   fmt.Sprintf("error updating task: %v", err),
		}
	}

	responseData, err := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: message,
	})
	if err != nil {
		return Response{
			Success: false,
			Error:   fmt.Sprintf("error marshaling response: %v", err),
		}
	}

	return Response{
		Success: true,
		Data:    responseData,
	}
}
