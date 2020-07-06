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

	flag.StringVar(&sqlitePath, "sqlite", "sqlite3://local.db", "path to sqlite file")
	flag.StringVar(&logPath, "log", "", "path to log file")
	flag.Parse()

	querier, err := db.Connect(sqlitePath)
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

	err = converter.ConvertFile(logFile, querier)
	if err != nil {
		fmt.Println("Error while parsing log file : ", err)
	}

}
