package svc

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/adrianpk/tyn/internal/model"
)

type ParserFunc func(input string) (model.Node, error)

type Svc struct {
	Repo   Repo
	Parser ParserFunc
}

func New(repo Repo) *Svc {
	return &Svc{
		Repo:   repo,
		Parser: Parse,
	}
}

// Capture parses the input text and stores the resulting node
func (s *Svc) Capture(text string) (model.Node, error) {
	node, err := s.Parser(text)
	if err != nil {
		return model.Node{}, err
	}

	node.GenID()

	err = s.Repo.Create(context.Background(), node)
	if err != nil {
		return model.Node{}, err
	}

	return node, nil
}

func (s *Svc) List(filter model.Filter) ([]model.Node, error) {
	nodes, err := s.Repo.List(context.Background())
	if err != nil {
		return nil, err
	}

	if isEmptyFilter(filter) {
		return nodes, nil
	}

	var filteredNodes []model.Node
	for _, node := range nodes {
		if matches(node, filter) {
			filteredNodes = append(filteredNodes, node)
		}
	}

	return filteredNodes, nil
}

func isEmptyFilter(filter model.Filter) bool {
	return filter.Type == "" && filter.Status == "" && len(filter.Tags) == 0 && len(filter.Places) == 0
}

func matches(node model.Node, filter model.Filter) bool {
	if filter.Type != "" && node.Type != filter.Type {
		return false
	}

	if filter.Status != "" && node.Status != filter.Status {
		return false
	}

	if len(filter.Tags) > 0 {
		tagMatch := false
		for _, tag := range filter.Tags {
			if hasTag(node, tag) {
				tagMatch = true
				break
			}
		}
		if !tagMatch {
			return false
		}
	}

	if len(filter.Places) > 0 {
		placeMatch := false
		for _, place := range filter.Places {
			if hasPlace(node, place) {
				placeMatch = true
				break
			}
		}
		if !placeMatch {
			return false
		}
	}

	return true
}

func hasTag(node model.Node, tag string) bool {
	for _, t := range node.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func hasPlace(node model.Node, place string) bool {
	for _, p := range node.Places {
		if p == place {
			return true
		}
	}
	return false
}

// GetOverdueTasks retrieves tasks with due dates in the past that aren't marked as done
// and haven't been notified today
func (s *Svc) GetOverdueTasks(ctx context.Context) ([]model.Node, error) {
	return s.Repo.GetOverdueTasks(ctx, model.NotificationType.DueDate)
}

// NotifyOverdueTask creates or updates a notification record for an overdue task
// and returns whether the notification was created (true) or updated (false)
func (s *Svc) NotifyOverdueTask(ctx context.Context, nodeID string) (bool, error) {
	// Check if notification already exists for this node
	notification, err := s.Repo.GetNotificationByNodeAndType(ctx, nodeID, model.NotificationType.DueDate)
	if err != nil {
		if err == sql.ErrNoRows {
			notification = model.Notification{
				NodeID:           nodeID,
				NotificationType: model.NotificationType.DueDate,
				LastNotifiedAt:   time.Now(),
				TimesNotified:    1,
			}
			notification.GenID()
			err = s.Repo.CreateNotification(ctx, notification)
			if err != nil {
				return false, err
			}
			return true, nil
		}
		return false, err
	}

	err = s.Repo.UpdateNotification(ctx, notification.ID, time.Now())
	if err != nil {
		return false, err
	}
	return false, nil
}

// GetNotificationByNodeAndType retrieves a notification by node id and notification type
func (s *Svc) GetNotificationByNodeAndType(ctx context.Context, nodeID, notificationType string) (model.Notification, error) {
	return s.Repo.GetNotificationByNodeAndType(ctx, nodeID, notificationType)
}

// GetAllTasks retrieves all tasks from the repository regardless of their creation date
func (s *Svc) GetAllTasks(ctx context.Context) ([]model.Node, error) {
	return s.Repo.GetAllTasks(ctx)
}

func (s *Svc) UpdateTask(ctx context.Context, id string, tags, places []string, dueDate string, text string) error {
	task, err := s.Repo.GetTaskByID(ctx, id)
	if err != nil {
		return fmt.Errorf("error retrieving task: %w", err)
	}

	if tags != nil {
		task.Tags = tags
	}

	if places != nil {
		task.Places = places
	}

	if dueDate != "" {
		date, err := time.Parse("2006-01-02", dueDate)
		if err != nil {
			return fmt.Errorf("invalid due date format: %w", err)
		}
		task.DueDate = &date
	}

	if text != "" {
		task.Content = text
	}

	err = s.Repo.UpdateTask(ctx, task)
	if err != nil {
		return fmt.Errorf("error updating task: %w", err)
	}

	return nil
}
