package svc

import (
	"context"
	"time"

	"github.com/adrianpk/tyn/internal/model"
)

// Repo defines the repository interface needed by the service
type Repo interface {
	// Node
	Create(ctx context.Context, node model.Node) error
	Get(ctx context.Context, id string) (model.Node, error)
	Update(ctx context.Context, node model.Node) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]model.Node, error)
	GetNodesByDay(day time.Time) ([]model.Node, error)
	// Notification
	CreateNotification(ctx context.Context, notification model.Notification) error
	GetNotification(ctx context.Context, id string) (model.Notification, error)
	GetNotificationByNodeAndType(ctx context.Context, nodeID, notificationType string) (model.Notification, error)
	UpdateNotification(ctx context.Context, id string, lastNotifiedAt time.Time) error
	DeleteNotification(ctx context.Context, id string) error
	DeleteNotificationByNode(ctx context.Context, nodeID string) error
	ListNotifications(ctx context.Context) ([]model.Notification, error)
	GetOverdueTasks(ctx context.Context, notificationType string) ([]model.Node, error)
}
