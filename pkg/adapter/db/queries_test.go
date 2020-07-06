package db

import (
	"app/pkg/abstract"
	"context"
	"os"
	"testing"
)

func TestCanInsertInDb(t *testing.T) {
	dbpath := "queries_test.db"

	defer os.Remove(dbpath)

	querier, err := Connect("sqlite3://" + dbpath)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer querier.Close()

	err = querier.Insert(context.TODO(), abstract.Line{})
	if err != nil {
		t.Fatalf("%v", err)
	}

}
