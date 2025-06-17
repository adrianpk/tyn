package sqlite

var Query = map[string]string{
	"create_table": `CREATE TABLE IF NOT EXISTS nodes (
		id TEXT PRIMARY KEY,
		type TEXT,
		content TEXT,
		link TEXT,
		tags TEXT,
		places TEXT,
		status TEXT,
		date DATETIME,
		override_date DATETIME
	);`,
	"create": `INSERT INTO nodes (id, type, content, link, tags, places, status, date, override_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
	"get":    `SELECT id, type, content, link, tags, places, status, date, override_date FROM nodes WHERE id = ?`,
	"update": `UPDATE nodes SET type=?, content=?, link=?, tags=?, places=?, status=?, date=?, override_date=? WHERE id=?`,
	"delete": `DELETE FROM nodes WHERE id = ?`,
	"list":   `SELECT id, type, content, link, tags, places, status, date, override_date FROM nodes`,
}
