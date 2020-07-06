package db

import (
	"app/pkg/abstract"
	"context"
	"database/sql"
)

type querier struct {
	db *sql.DB
}

func (q *querier) Insert(ctx context.Context, line abstract.Line) error {
	return nil
}

func (q *querier) Close() error {
	return q.db.Close()
}
