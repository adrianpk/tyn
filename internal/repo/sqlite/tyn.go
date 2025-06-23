package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrianpk/tyn/internal/config"
	"github.com/adrianpk/tyn/internal/model"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type TynRepo struct {
	db  *sqlx.DB
	cfg *config.Config
}

func NewTynRepo(cfg *config.Config) (*TynRepo, error) {
	dbPath := getDBPath()
	db, err := sqlx.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	err = migrate(db)
	if err != nil {
		return nil, err
	}

	return &TynRepo{db: db, cfg: cfg}, nil
}

func getDBPath() string {
	if dbPath := os.Getenv("TYN_DB_PATH"); dbPath != "" {
		log.Printf("Using database path from environment: %s", dbPath)
		return dbPath
	}

	if os.Getenv("TYN_DEV") != "" {
		log.Printf("Development mode: using local tyn.db")
		return "tyn.db"
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Warning: could not get user home directory: %v", err)
		return "tyn.db"
	}

	tynDir := filepath.Join(homeDir, ".tyn")

	if err := os.MkdirAll(tynDir, 0755); err != nil {
		log.Printf("Warning: could not create directory %s: %v", tynDir, err)
		return "tyn.db"
	}

	dbPath := filepath.Join(tynDir, "tyn.db")
	log.Printf("Using database at: %s", dbPath)
	return dbPath
}

func (r *TynRepo) Create(ctx context.Context, node model.Node) error {
	log.Printf("Repository - Before writing to DB: DueDate = %v", node.DueDate)

	var dueDateStr interface{} = nil
	if node.DueDate != nil {
		utcDueDate := node.DueDate.UTC()
		dueDateStr = utcDueDate.Format(model.DateTimeFormat)
		log.Printf("Repository - Formatted DueDate for DB (UTC): %v", dueDateStr)
	}

	_, err := r.db.ExecContext(ctx, Query["create"],
		node.ID, node.Type, node.Content, node.Link,
		stringSliceToCSV(node.Tags), stringSliceToCSV(node.Places), node.Status,
		node.Draft, node.Date.UTC().Format(model.DateTimeFormat), dueDateStr,
	)
	return err
}

func (r *TynRepo) Get(ctx context.Context, id string) (model.Node, error) {
	var node model.Node
	var tags, places string
	var dueDate sql.NullTime

	row := r.db.QueryRowContext(ctx, Query["get"], id)
	err := row.Scan(&node.ID, &node.Type, &node.Content, &node.Link, &tags, &places, &node.Status, &node.Draft, &node.Date, &dueDate)
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
		dueDateStr = utcDueDate.Format(model.DateTimeFormat)
	}

	_, err := r.db.ExecContext(ctx, Query["update"],
		node.Type, node.Content, node.Link,
		stringSliceToCSV(node.Tags), stringSliceToCSV(node.Places), node.Status,
		node.Draft, node.Date.UTC().Format(model.DateTimeFormat), dueDateStr, node.ID,
	)
	return err
}

func (r *TynRepo) UpdateTask(ctx context.Context, node model.Node) error {
	var dueDateStr interface{} = nil
	if node.DueDate != nil {
		utcDueDate := node.DueDate.UTC()
		dueDateStr = utcDueDate.Format(model.DateTimeFormat)
	}

	_, err := r.db.ExecContext(ctx, Query["update"],
		node.Type, node.Content, node.Link,
		stringSliceToCSV(node.Tags), stringSliceToCSV(node.Places), node.Status,
		node.Draft, node.Date.UTC().Format(model.DateTimeFormat), dueDateStr, node.ID,
	)
	return err
}

func (r *TynRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, Query["delete"], id)
	return err
}

func (r *TynRepo) List(ctx context.Context) ([]model.Node, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	daysLimit := r.cfg.DoneTaskListDays
	cutoff := time.Now().AddDate(0, 0, -daysLimit)
	cutoffStr := cutoff.Format(model.DateTimeFormat)

	query := `SELECT id, type, content, link, tags, places, status, draft, date, due_date FROM nodes
		WHERE type != 'task'
		   OR (
			   type = 'task' AND (
				   status NOT IN ('done', 'canceled')
				   OR (status IN ('done', 'canceled') AND date >= ?)
			   )
		   )`

	rows, err := r.db.QueryContext(ctx, query, cutoffStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []model.Node
	for rows.Next() {
		var node model.Node
		var tags, places string
		var dueDate sql.NullTime

		err := rows.Scan(&node.ID, &node.Type, &node.Content, &node.Link, &tags, &places, &node.Status, &node.Draft, &node.Date, &dueDate)
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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return nodes, nil
}

func (r *TynRepo) GetNodesByDay(day time.Time) ([]model.Node, error) {
	ctx := context.Background()

	startLocal := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.Local)
	endLocal := startLocal.Add(24 * time.Hour)
	startUTC := startLocal.UTC()
	endUTC := endLocal.UTC()

	rows, err := r.db.QueryContext(ctx, Query["list_by_day"],
		startUTC.Format(model.DateTimeFormat),
		endUTC.Format(model.DateTimeFormat))
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
			&node.ID, &node.Type, &node.Content, &node.Link,
			&tags, &places, &node.Status, &node.Draft, &node.Date, &dueDate,
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

	log.Printf("GetNodesByDay - Found %d nodes for date %s", len(nodes), day.Format("2006-01-02"))
	return nodes, nil
}

func (r *TynRepo) GetNotesAndLinksByDay(day time.Time) ([]model.Node, error) {
	ctx := context.Background()

	startLocal := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.Local)
	endLocal := startLocal.Add(24 * time.Hour)
	startUTC := startLocal.UTC()
	endUTC := endLocal.UTC()

	rows, err := r.db.QueryContext(ctx, Query["list_notes_and_links_by_day"],
		startUTC.Format(model.DateTimeFormat),
		endUTC.Format(model.DateTimeFormat))
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
			&node.ID, &node.Type, &node.Content, &node.Link,
			&tags, &places, &node.Status, &node.Draft, &node.Date, &dueDate,
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

	log.Printf("GetNotesAndLinksByDay - Found %d notes and links for date %s", len(nodes), day.Format("2006-01-02"))
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

func (r *TynRepo) GetOverdueTasks(ctx context.Context, notificationType string) ([]model.Node, error) {
	log.Printf("Executing get_overdue_tasks with notificationType: %s", notificationType)
	query := Query["get_overdue_tasks"]

	now := time.Now()
	startLocal := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	endLocal := startLocal.Add(24 * time.Hour)
	startUTC := startLocal.UTC()
	endUTC := endLocal.UTC()

	rows, err := r.db.QueryContext(ctx, query, notificationType,
		startUTC.Format(model.DateTimeFormat),
		endUTC.Format(model.DateTimeFormat))
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
			&node.ID, &node.Type, &node.Content, &node.Link,
			&tags, &places, &node.Status, &node.Draft, &node.Date, &dueDate,
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

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return nodes, nil
}

func (r *TynRepo) GetAllTasks(ctx context.Context) ([]model.Node, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	rows, err := r.db.QueryContext(ctx, Query["list_all_tasks"])
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Node
	for rows.Next() {
		var task model.Node
		var tags, places string
		var dueDate sql.NullTime

		err := rows.Scan(
			&task.ID, &task.Type, &task.Content, &task.Link,
			&tags, &places, &task.Status, &task.Draft, &task.Date, &dueDate,
		)
		if err != nil {
			return nil, err
		}

		task.Tags = csvToStringSlice(tags)
		task.Places = csvToStringSlice(places)
		if dueDate.Valid {
			localTime := dueDate.Time.In(time.Local)
			task.DueDate = &localTime
		}

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TynRepo) GetTaskByID(ctx context.Context, id string) (model.Node, error) {
	var node model.Node
	var tags, places string
	var dueDate sql.NullTime

	row := r.db.QueryRowContext(ctx, Query["get"], id)
	err := row.Scan(&node.ID, &node.Type, &node.Content, &node.Link, &tags, &places, &node.Status, &node.Draft, &node.Date, &dueDate)
	if err == nil {
		node.Tags = csvToStringSlice(tags)
		node.Places = csvToStringSlice(places)
		if dueDate.Valid {
			localTime := dueDate.Time.In(time.Local)
			node.DueDate = &localTime
		}
		return node, nil
	}

	rows, err := r.db.QueryContext(ctx, Query["get_by_partial_id"], id)
	if err != nil {
		return model.Node{}, fmt.Errorf("task with ID '%s' not found: %v", id, err)
	}
	defer rows.Close()

	var nodes []model.Node
	for rows.Next() {
		var n model.Node
		var t, p string
		var dd sql.NullTime

		err := rows.Scan(&n.ID, &n.Type, &n.Content, &n.Link, &t, &p, &n.Status, &n.Draft, &n.Date, &dd)
		if err != nil {
			return model.Node{}, fmt.Errorf("error scanning task: %v", err)
		}

		n.Tags = csvToStringSlice(t)
		n.Places = csvToStringSlice(p)
		if dd.Valid {
			localTime := dd.Time.In(time.Local)
			n.DueDate = &localTime
		}

		nodes = append(nodes, n)
	}

	if len(nodes) == 0 {
		return model.Node{}, fmt.Errorf("task with ID '%s' not found", id)
	}

	if len(nodes) > 1 {
		for _, n := range nodes {
			if strings.HasPrefix(n.ID, id) {
				return n, nil
			}
		}
		log.Printf("Warning: Multiple tasks found with ID prefix '%s', using first match", id)
	}

	return nodes[0], nil
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

func migrate(db *sqlx.DB) error {
	_, err := db.Exec(Query["create_nodes_table"])
	if err != nil {
		return err
	}

	_, err = db.Exec(Query["create_notifications_table"])
	if err != nil {
		return err
	}

	var hasDraft bool
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('nodes') WHERE name='draft'").Scan(&hasDraft)
	if err != nil {
		return err
	}

	if !hasDraft {
		log.Printf("Adding draft column to nodes table")
		_, err = db.Exec("ALTER TABLE nodes ADD COLUMN draft TEXT")
		if err != nil {
			return err
		}
	}

	return nil
}
