package main

import (
	"app/pkg/adapter/db"
	"app/pkg/domain/converter"
	"flag"
	"fmt"
	"os"
)

func main() {
	var sqlitePath string
	var logPath string
	var clearDb bool
	var batchSize int

	flag.StringVar(&sqlitePath, "sqlite", "local.db", "path to sqlite file")
	flag.StringVar(&logPath, "log", "", "path to log file")
	flag.BoolVar(&clearDb, "clearDb", false, "Should we clear database if already present ? ")
	flag.IntVar(&batchSize, "batchSize", 1000, "number of item to inster by query")
	flag.Parse()

	if clearDb {
		fmt.Println("Cleanning previous database")
		_ = os.Remove(sqlitePath)
	}

	querier, err := db.Connect("sqlite3://" + sqlitePath)
	if err != nil {
		fmt.Println("cannot connect", err)
		return
	}
	defer querier.Close()

	logFile, err := os.Open(logPath)
	if err != nil {
		fmt.Println("cant find log file : ", err)
		return
	}
	defer logFile.Close()

	fmt.Println("converting file : ", logPath)
	err = converter.ConvertFile(logFile, querier, batchSize)
	if err != nil {
		fmt.Println("Error while parsing log file : ", err)
	}

}
