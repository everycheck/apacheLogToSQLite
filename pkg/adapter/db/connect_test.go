package db

import (
	"testing"
)

func TestBadlyFormatedUrl(t *testing.T) {

	querier, err := Connect("sqlite3:/TestBadlyFormatedUrl")
	if err == nil {
		querier.Close()
		t.Fatalf("should have failed to connect")
	}
}

func TestWrongPath(t *testing.T) {

	querier, err := Connect("sqlite3:///TestBadlyFormatedUrl")
	if err == nil {
		querier.Close()
		t.Fatalf("should have failed to connect")
	}
}

func TestWrongProtocol(t *testing.T) {

	querier, err := Connect("sqlite://test.db")
	if err == nil {
		querier.Close()
		t.Fatalf("should have failed to connect")
	}
}

func TestCanConnectToExistingDB(t *testing.T) {
	dbpath := "connect_test.db"

	querier, err := Connect("sqlite3://" + dbpath)
	if err != nil {
		t.Fatalf("%v", err)
	}
	querier.Close()

}
