package bkg

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func (s *Service) handleDate(params json.RawMessage) Response {
	var dateParams DateCmdParams
	err := json.Unmarshal(params, &dateParams)
	if err != nil {
		return Response{
			Success: false,
			Error:   fmt.Sprintf("error unmarshaling date params: %v", err),
		}
	}

	log.Printf("Handling date operation: %s for task %s with date %s", dateParams.Operation, dateParams.ID, dateParams.Date)

	ctx := context.Background()
	task, err := s.svc.Repo.GetTaskByID(ctx, dateParams.ID)
	if err != nil {
		return Response{
			Success: false,
			Error:   fmt.Sprintf("error fetching task with ID %s: %v", dateParams.ID, err),
		}
	}

	var message string
	var originalDate string
	var newDate string

	if task.DueDate != nil {
		originalDate = task.DueDate.Format("2006-01-02")
	} else {
		originalDate = "none"
	}

	switch dateParams.Operation {
	case "set":
		date, err := time.Parse("2006-01-02", dateParams.Date)
		if err != nil {
			return Response{
				Success: false,
				Error:   fmt.Sprintf("invalid date format: %v", err),
			}
		}
		task.DueDate = &date
		newDate = dateParams.Date
		message = fmt.Sprintf("Set due date to %s for task %s", newDate, dateParams.ID)

	case "remove":
		task.DueDate = nil
		newDate = "none"
		message = fmt.Sprintf("Removed due date from task %s", dateParams.ID)

	default:
		return Response{
			Success: false,
			Error:   fmt.Sprintf("unknown date operation: %s", dateParams.Operation),
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
		Message      string `json:"message"`
		OriginalDate string `json:"original_date"`
		NewDate      string `json:"new_date,omitempty"`
	}{
		Message:      message,
		OriginalDate: originalDate,
		NewDate:      newDate,
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
