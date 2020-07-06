package db

import (
	"app/pkg/abstract"
	"context"
	"database/sql"
)

type querier struct {
	db *sql.DB
}

const insertQuery = `INSERT INTO entry (
	RemoteHost ,Time ,Request ,Status ,Bytes ,Referer ,UserAgent ,URL
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)`

func (q *querier) Insert(ctx context.Context, line abstract.Line) error {
	_, err := q.db.ExecContext(ctx, insertQuery,
		line.RemoteHost,
		line.Time,
		line.Request,
		line.Status,
		line.Bytes,
		line.Referer,
		line.UserAgent,
		line.URL,
	)
	return err
}

func (q *querier) Close() error {
	return q.db.Close()
}
