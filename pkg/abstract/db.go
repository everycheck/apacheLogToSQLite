package abstract

import (
	"context"
	"time"
)

type Line struct {
	RemoteHost string
	Time       time.Time
	Request    string
	Status     int
	Bytes      int
	Referer    string
	UserAgent  string
	URL        string
}

type DBLineInserter interface {
	Insert(ctx context.Context, line Line) error
}

type DBLineBulkInserter interface {
	BulkInsert(ctx context.Context, lines []Line) error
}

type DB interface {
	DBLineInserter
	DBLineBulkInserter
	Close() error
}
