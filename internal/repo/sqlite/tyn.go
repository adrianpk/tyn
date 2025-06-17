package sqlite

import (
	"context"
	"database/sql"

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
	_, err := r.db.ExecContext(ctx, Query["create"],
		node.ID, node.Type, node.Content, node.Link,
		stringSliceToCSV(node.Tags), stringSliceToCSV(node.Places), node.Status,
		node.Date, node.OverrideDate,
	)
	return err
}

func (r *TynRepo) Get(ctx context.Context, id string) (model.Node, error) {
	row := r.db.QueryRowContext(ctx, Query["get"], id)
	var node model.Node
	var tags, places string
	var overrideDate sql.NullTime
	if err := row.Scan(&node.ID, &node.Type, &node.Content, &node.Link, &tags, &places, &node.Status, &node.Date, &overrideDate); err != nil {
		return node, err
	}
	node.Tags = csvToStringSlice(tags)
	node.Places = csvToStringSlice(places)
	if overrideDate.Valid {
		node.OverrideDate = &overrideDate.Time
	}
	return node, nil
}

func (r *TynRepo) Update(ctx context.Context, node model.Node) error {
	_, err := r.db.ExecContext(ctx, Query["update"],
		node.Type, node.Content, node.Link,
		stringSliceToCSV(node.Tags), stringSliceToCSV(node.Places), node.Status,
		node.Date, node.OverrideDate, node.ID,
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
		var overrideDate sql.NullTime
		if err := rows.Scan(&node.ID, &node.Type, &node.Content, &node.Link, &tags, &places, &node.Status, &node.Date, &overrideDate); err != nil {
			return nil, err
		}
		node.Tags = csvToStringSlice(tags)
		node.Places = csvToStringSlice(places)
		if overrideDate.Valid {
			node.OverrideDate = &overrideDate.Time
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func stringSliceToCSV(s []string) string {
	return model.EncodeCSV(s)
}

func csvToStringSlice(s string) []string {
	return model.DecodeCSV(s)
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(Query["create_table"])
	return err
}
