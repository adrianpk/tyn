package sqlite

var Query = map[string]string{
	"create_nodes_table": `CREATE TABLE IF NOT EXISTS nodes (
		id TEXT PRIMARY KEY,
		type TEXT,
		content TEXT,
		link TEXT,
		tags TEXT,
		places TEXT,
		status TEXT,
		draft TEXT,
		date DATETIME,
		due_date DATETIME
	);`,
	"create_notifications_table": `CREATE TABLE IF NOT EXISTS notifications (
		id TEXT PRIMARY KEY,
		node_id TEXT NOT NULL,
		notification_type TEXT NOT NULL,
		last_notified_at TIMESTAMP NOT NULL,
		times_notified INTEGER DEFAULT 1,
		FOREIGN KEY (node_id) REFERENCES nodes (id) ON DELETE CASCADE
	);`,

	// Node queries
	"create":            `INSERT INTO nodes (id, type, content, link, tags, places, status, draft, date, due_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
	"get":               `SELECT id, type, content, link, tags, places, status, draft, date, due_date FROM nodes WHERE id = ?`,
	"get_by_partial_id": `SELECT id, type, content, link, tags, places, status, draft, date, due_date FROM nodes WHERE id LIKE ? || '%'`,
	"update":            `UPDATE nodes SET type=?, content=?, link=?, tags=?, places=?, status=?, draft=?, date=?, due_date=? WHERE id=?`,
	"delete":            `DELETE FROM nodes WHERE id = ?`,
	"list":              `SELECT id, type, content, link, tags, places, status, draft, date, due_date FROM nodes`,
	"list_by_day": `SELECT id, type, content, link, tags, places, status, draft, date, due_date FROM nodes 
		WHERE date >= ? AND date < ?`,
	"list_notes_and_links_by_day": `SELECT id, type, content, link, tags, places, status, draft, date, due_date FROM nodes 
		WHERE (type = 'note' OR type = 'link') AND date >= ? AND date < ?`,
	"list_all_tasks": `SELECT id, type, content, link, tags, places, status, draft, date, due_date FROM nodes
		WHERE type = 'task' ORDER BY date`,

	// Notification queries
	"create_notification": `INSERT INTO notifications (id, node_id, notification_type, last_notified_at, times_notified) 
		VALUES (?, ?, ?, ?, ?)`,
	"get_notification": `SELECT id, node_id, notification_type, last_notified_at, times_notified 
		FROM notifications 
		WHERE id = ?`,
	"get_notification_by_node_and_type": `SELECT id, node_id, notification_type, last_notified_at, times_notified 
		FROM notifications 
		WHERE node_id = ? AND notification_type = ?`,
	"update_notification": `UPDATE notifications 
		SET last_notified_at = ?, times_notified = times_notified + 1 
		WHERE id = ?`,
	"delete_notification":         `DELETE FROM notifications WHERE id = ?`,
	"delete_notification_by_node": `DELETE FROM notifications WHERE node_id = ?`,
	"list_notifications": `SELECT id, node_id, notification_type, last_notified_at, times_notified 
		FROM notifications`,
	"get_overdue_tasks": `SELECT n.id, n.type, n.content, n.link, n.tags, n.places, n.status, n.draft, n.date, n.due_date
		FROM nodes n
		WHERE n.type = 'task'
		AND n.due_date IS NOT NULL
		AND strftime('%s', n.due_date) < strftime('%s', 'now')
		AND (n.status IS NULL OR n.status != 'done')
		AND NOT EXISTS (
			SELECT 1 FROM notifications nt
			WHERE nt.node_id = n.id
			AND nt.notification_type = ?
			AND nt.last_notified_at >= ? AND nt.last_notified_at < ?
		)`,
}
