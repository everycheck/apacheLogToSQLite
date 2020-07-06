package db

import (
	"os"
	"testing"
)

func TestCanInitDb(t *testing.T) {
	dbpath := "int_test.db"

	defer os.Remove(dbpath)

	querier, err := Connect("sqlite3://" + dbpath)
	if err != nil {
		t.Fatalf("%v", err)
	}
	querier.Close()
}
