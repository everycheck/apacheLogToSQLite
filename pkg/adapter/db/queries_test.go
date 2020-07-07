package db

import (
	"app/pkg/abstract"
	"context"
	"os"
	"testing"
)

func TestbuildValuesQuery(t *testing.T) {
	tests := []struct {
		nbField  int
		nbItem   int
		expected string
	}{
		{1, 1, "(%1)"},
		{3, 1, "(%1, %2, %3)"},
		{3, 2, "(%1, %2, %3), (%4, %5, %6)"},
		{-1, 2, ""},
		{3, -1, "(), ()"},
	}

	for i, tt := range tests {
		result := buildValuesQuery(tt.nbField, tt.nbItem)
		if result != tt.expected {
			t.Fatalf("[%d] buildValuesQuery failed exp: %s, got: %s\n", i, tt.expected, result)
		}
	}
}

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

func TestCanBulkInsertInDb(t *testing.T) {
	dbpath := "queries_test.db"

	defer os.Remove(dbpath)

	querier, err := Connect("sqlite3://" + dbpath)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer querier.Close()

	err = querier.BulkInsert(context.TODO(), []abstract.Line{{}, {}, {}})
	if err != nil {
		t.Fatalf("%v", err)
	}
}
