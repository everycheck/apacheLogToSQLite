package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func wasDbAlreadyCreated(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}

	return true
}

func Connect(url string) (*querier, error) {
	urlPart := strings.Split(url, "://")

	if len(urlPart) != 2 {
		return nil, fmt.Errorf("Malformated DB url : expect {protocol}://{detail}")
	}

	switch urlPart[0] {
	case "sqlite3":
		wasDbCreated := wasDbAlreadyCreated(urlPart[1])
		db, err := sql.Open(urlPart[0], urlPart[1])
		if err != nil {
			return nil, fmt.Errorf("Cannot open %s : %w", urlPart[1], err)
		}

		if !wasDbCreated {
			err = initTable(context.TODO(), db)
			if err != nil {
				return nil, fmt.Errorf("Cannot initialize db at %s : %w", urlPart[1], err)
			}
		}

		return &querier{db: db}, nil
	}

	return nil, fmt.Errorf("Unrecognized protocol :%s, expect : sqlite3\n", urlPart[0])
}
