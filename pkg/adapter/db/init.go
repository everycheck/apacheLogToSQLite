package db

import (
	"context"
	"database/sql"
	"fmt"
)

const initQuery = `DROP TABLE IF EXISTS entry;
CREATE TABLE entry (
	RemoteHost CHAR(255)
	Time       DATETIME
	Request    CHAR(255)
	Status     INTEGER
	Bytes      INTEGER
	Referer    CHAR(255)
	UserAgent  CHAR(255)
	URL        CHAR(255)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;`

func initTable(ctx context.Context, db *sql.DB) error {
	fmt.Println("Local database initialization...")
	_, err := db.ExecContext(ctx, initQuery)
	return err
}
