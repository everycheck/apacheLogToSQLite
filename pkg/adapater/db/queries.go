package db

import (
	"app/pkg/abstract"
)

type db struct {
	db *sql.DB
}

func (d *db) Insert(ctx context.context, line abstract.Line) error {
	return nil
}

func (d *db) Close() {
	return d.db.Close()
}
