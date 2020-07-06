package db

import (
	"context"
	"os"
	"testing"
)

func TestCanInitDb(t *testing.T) {
	dbpath := "init_test.db"

	defer os.Remove(dbpath)

	querier, err := Connect("sqlite3://" + dbpath)
	if err != nil {
		t.Fatalf("%v", err)
	}
	querier.Close()
}

func TestCanInitDbCalledTwiceShouldFailed(t *testing.T) {
	dbpath := "init_test.db"

	defer os.Remove(dbpath)

	querier, err := Connect("sqlite3://" + dbpath)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer querier.Close()

	err = initTable(context.TODO(), querier.db)
	if err == nil {
		t.Fatalf("Should have failed when intializing twice")
	}
}
