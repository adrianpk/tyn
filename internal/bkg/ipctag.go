package bkg

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/adrianpk/tyn/internal/model"
)

func (s *Service) handleTag(params json.RawMessage) Response {
	var tagParams TagCmdParams
	err := json.Unmarshal(params, &tagParams)
	if err != nil {
		return Response{
			Success: false,
			Error:   fmt.Sprintf("error unmarshaling tag params: %v", err),
		}
	}

	log.Printf("Handling tag operation: %s for task %s with tags %v", tagParams.Operation, tagParams.ID, tagParams.Tags)

	ctx := context.Background()
	task, err := s.svc.Repo.GetTaskByID(ctx, tagParams.ID)
	if err != nil {
		return Response{
			Success: false,
			Error:   fmt.Sprintf("error fetching task with ID %s: %v", tagParams.ID, err),
		}
	}

	var message string

	switch tagParams.Operation {
	case "add":
		for _, tag := range tagParams.Tags {
			exists := false
			for _, existingTag := range task.Tags {
				if existingTag == tag {
					exists = true
					break
				}
			}
			if !exists {
				task.Tags = append(task.Tags, tag)
			}
		}
		message = fmt.Sprintf("Added tags %v to task %s", tagParams.Tags, tagParams.ID)

	case "remove":
		for _, tagToRemove := range tagParams.Tags {
			newTags := []string{}
			for _, existingTag := range task.Tags {
				if existingTag != tagToRemove {
					newTags = append(newTags, existingTag)
				}
			}
			task.Tags = newTags
		}
		message = fmt.Sprintf("Removed tags %v from task %s", tagParams.Tags, tagParams.ID)

	case "clear":
		task.Tags = []string{}
		message = fmt.Sprintf("Cleared all tags from task %s", tagParams.ID)

	default:
		return Response{
			Success: false,
			Error:   fmt.Sprintf("unknown tag operation: %s", tagParams.Operation),
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

func containsString(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

func (s *Service) findTaskByID(ctx context.Context, id string) (model.Node, error) {
	nodes, err := s.svc.Repo.List(ctx)
	if err != nil {
		log.Printf("Error listing nodes: %v", err)
		return model.Node{}, fmt.Errorf("error listing nodes: %v", err)
	}

	var targetNode model.Node
	var found bool

	for _, node := range nodes {
		if len(id) <= 4 && node.ShortID() == id {
			targetNode = node
			found = true
			log.Printf("Found node by ShortID: %s", node.ID)
			break
		} else if len(id) > 4 && len(node.ID) >= len(id) && node.ID[:len(id)] == id {
			targetNode = node
			found = true
			log.Printf("Found node by prefix: %s", node.ID)
			break
		}
	}

	if !found {
		log.Printf("Task with ID '%s' not found", id)
		return model.Node{}, fmt.Errorf("task with ID '%s' not found", id)
	}

	if targetNode.Type != "task" {
		log.Printf("Node with ID '%s' is not a task", id)
		return model.Node{}, fmt.Errorf("node with ID '%s' is not a task", id)
	}

	return targetNode, nil
}
