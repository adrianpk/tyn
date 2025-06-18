package sqlite

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"github.com/adrianpk/tyn/internal/model"
	_ "modernc.org/sqlite"
)

type TynRepo struct {
	db *sql.DB
}

func NewTynRepo() (*TynRepo, error) {
	db, err := sql.Open("sqlite", "tyn.db")
	if err != nil {
		return nil, err
	}

	err = migrate(db)
	if err != nil {
		return nil, err
	}
	return &TynRepo{db: db}, nil
}

func (r *TynRepo) Create(ctx context.Context, node model.Node) error {
	log.Printf("Repository - Before writing to DB: DueDate = %v", node.DueDate)

	var dueDateStr interface{} = nil
	if node.DueDate != nil {
		utcDueDate := node.DueDate.UTC()
		dueDateStr = utcDueDate.Format("2006-01-02 15:04:05")
		log.Printf("Repository - Formatted DueDate for DB (UTC): %v", dueDateStr)
	}

	_, err := r.db.ExecContext(ctx, Query["create"],
		node.ID, node.Type, node.Content, node.Link,
		stringSliceToCSV(node.Tags), stringSliceToCSV(node.Places), node.Status,
		node.Date.UTC().Format("2006-01-02 15:04:05"), dueDateStr,
	)
	return err
}

func (r *TynRepo) Get(ctx context.Context, id string) (model.Node, error) {
	row := r.db.QueryRowContext(ctx, Query["get"], id)
	var node model.Node
	var tags, places string
	var dueDate sql.NullTime
	err := row.Scan(&node.ID, &node.Type, &node.Content, &node.Link, &tags, &places, &node.Status, &node.Date, &dueDate)
	if err != nil {
		return node, err
	}
	node.Tags = csvToStringSlice(tags)
	node.Places = csvToStringSlice(places)
	if dueDate.Valid {
		localTime := dueDate.Time.In(time.Local)
		node.DueDate = &localTime
	}
	return node, nil
}

func (r *TynRepo) Update(ctx context.Context, node model.Node) error {
	var dueDateStr interface{} = nil
	if node.DueDate != nil {
		utcDueDate := node.DueDate.UTC()
		dueDateStr = utcDueDate.Format("2006-01-02 15:04:05")
	}

	_, err := r.db.ExecContext(ctx, Query["update"],
		node.Type, node.Content, node.Link,
		stringSliceToCSV(node.Tags), stringSliceToCSV(node.Places), node.Status,
		node.Date.UTC().Format("2006-01-02 15:04:05"), dueDateStr, node.ID,
	)
	return err
}

func (r *TynRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, Query["delete"], id)
	return err
}

func (r *TynRepo) List(ctx context.Context) ([]model.Node, error) {
	rows, err := r.db.QueryContext(ctx, Query["list"])
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var nodes []model.Node
	for rows.Next() {
		var node model.Node
		var tags, places string
		var dueDate sql.NullTime
		err := rows.Scan(&node.ID, &node.Type, &node.Content, &node.Link, &tags, &places, &node.Status, &node.Date, &dueDate)
		if err != nil {
			return nil, err
		}

		log.Printf("Repository - After DB read: Raw dueDate = %v, Valid = %v", dueDate.Time, dueDate.Valid)

		node.Tags = csvToStringSlice(tags)
		node.Places = csvToStringSlice(places)
		if dueDate.Valid {
			localTime := dueDate.Time.In(time.Local)
			node.DueDate = &localTime
			log.Printf("Repository - After setting pointer: node.DueDate = %v", *node.DueDate)
		}
		nodes = append(nodes, node)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (r *TynRepo) GetNodesByDay(day time.Time) ([]model.Node, error) {
	ctx := context.Background()

	dateStr := day.Format("2006-01-02")

	rows, err := r.db.QueryContext(ctx, Query["list_by_day"], dateStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []model.Node
	for rows.Next() {
		var node model.Node
		var tags, places string
		var dueDate sql.NullTime

		err := rows.Scan(
			&node.ID,
			&node.Type,
			&node.Content,
			&node.Link,
			&tags,
			&places,
			&node.Status,
			&node.Date,
			&dueDate,
		)
		if err != nil {
			return nil, err
		}

		node.Tags = csvToStringSlice(tags)
		node.Places = csvToStringSlice(places)
		if dueDate.Valid {
			localTime := dueDate.Time.In(time.Local)
			node.DueDate = &localTime
		}

		nodes = append(nodes, node)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (r *TynRepo) CreateNotification(ctx context.Context, notification model.Notification) error {
	_, err := r.db.ExecContext(ctx, Query["create_notification"],
		notification.ID,
		notification.NodeID,
		notification.NotificationType,
		notification.LastNotifiedAt,
		notification.TimesNotified,
	)
	return err
}

func (r *TynRepo) GetNotification(ctx context.Context, id string) (model.Notification, error) {
	row := r.db.QueryRowContext(ctx, Query["get_notification"], id)
	var notification model.Notification
	err := row.Scan(
		&notification.ID,
		&notification.NodeID,
		&notification.NotificationType,
		&notification.LastNotifiedAt,
		&notification.TimesNotified,
	)
	if err != nil {
		return notification, err
	}
	return notification, nil
}

func (r *TynRepo) GetNotificationByNodeAndType(ctx context.Context, nodeID, notificationType string) (model.Notification, error) {
	row := r.db.QueryRowContext(ctx, Query["get_notification_by_node_and_type"], nodeID, notificationType)
	var notification model.Notification
	err := row.Scan(
		&notification.ID,
		&notification.NodeID,
		&notification.NotificationType,
		&notification.LastNotifiedAt,
		&notification.TimesNotified,
	)

	if err != nil {
		return notification, err
	}
	return notification, nil
}

func (r *TynRepo) UpdateNotification(ctx context.Context, id string, lastNotifiedAt time.Time) error {
	_, err := r.db.ExecContext(ctx, Query["update_notification"], lastNotifiedAt, id)
	return err
}

func (r *TynRepo) DeleteNotification(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, Query["delete_notification"], id)
	return err
}

func (r *TynRepo) DeleteNotificationByNode(ctx context.Context, nodeID string) error {
	_, err := r.db.ExecContext(ctx, Query["delete_notification_by_node"], nodeID)
	return err
}

func (r *TynRepo) ListNotifications(ctx context.Context) ([]model.Notification, error) {
	rows, err := r.db.QueryContext(ctx, Query["list_notifications"])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []model.Notification
	for rows.Next() {
		var notification model.Notification
		err := rows.Scan(
			&notification.ID,
			&notification.NodeID,
			&notification.NotificationType,
			&notification.LastNotifiedAt,
			&notification.TimesNotified,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

// GetOverdueTasks retrieves tasks with due dates in the past that aren't marked as done
// and haven't been notified today for the specified notification type
func (r *TynRepo) GetOverdueTasks(ctx context.Context, notificationType string) ([]model.Node, error) {
	log.Printf("Executing get_overdue_tasks with notificationType: %s", notificationType)
	query := Query["get_overdue_tasks"]
	rows, err := r.db.QueryContext(ctx, query, notificationType)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var nodes []model.Node
	for rows.Next() {
		var node model.Node
		var tags, places string
		var dueDate sql.NullTime

		err := rows.Scan(
			&node.ID,
			&node.Type,
			&node.Content,
			&node.Link,
			&tags,
			&places,
			&node.Status,
			&node.Date,
			&dueDate,
		)
		if err != nil {
			return nil, err
		}

		node.Tags = csvToStringSlice(tags)
		node.Places = csvToStringSlice(places)
		if dueDate.Valid {
			localTime := dueDate.Time.In(time.Local)
			node.DueDate = &localTime
			log.Printf("Found overdue task with due date: %v", *node.DueDate)
		}

		nodes = append(nodes, node)
	}

	log.Printf("Found %d overdue tasks", len(nodes))

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func stringSliceToCSV(s []string) string {
	return strings.Join(s, ",")
}

func csvToStringSlice(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(Query["create_nodes_table"])
	if err != nil {
		return err
	}

	_, err = db.Exec(Query["create_notifications_table"])
	if err != nil {
		return err
	}

	return nil
}
