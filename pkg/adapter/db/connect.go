package db

import (
	"app/pkg/abstract"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func Connect(url string) (abstract.DB, error) {
	urlPart := strings.Split(url, "://")

	if len(urlPart) != 2 {
		return nil, fmt.Errorf("Malformated DB url : expect {protocol}://{detail}")
	}

	switch urlPart[0] {
	case "sqlite3":
		db, err := sql.Open(urlPart[0], urlPart[1])
		if err != nil {
			return nil, fmt.Errorf("Cannot open %s : %w", urlPart[1], err)
		}
		return &querier{db: db}, nil
	}

	return nil, fmt.Errorf("Unrecognized protocol :%s, expect : sqlite3\n", urlPart[0])
}
