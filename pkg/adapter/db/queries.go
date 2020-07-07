package db

import (
	"app/pkg/abstract"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type querier struct {
	db         *sql.DB
	insertStmt map[int]*sql.Stmt
}

const insertQuery = `INSERT INTO entry (
	RemoteHost ,Time ,Request ,Status ,Bytes ,Referer ,UserAgent ,URL
) VALUES %s`

func (q *querier) Insert(ctx context.Context, line abstract.Line) error {
	return q.BulkInsert(ctx, []abstract.Line{line})
}

func (q *querier) BulkInsert(ctx context.Context, lines []abstract.Line) error {
	nbFieldInsterted := 8
	var err error
	stmt, ok := q.insertStmt[len(lines)]
	if !ok {
		query := fmt.Sprintf(insertQuery, buildValuesQuery(nbFieldInsterted, len(lines)))
		stmt, err = q.db.Prepare(query)
		if err != nil {
			return fmt.Errorf("Error while preparing insert statment : %w", err)
		}
		q.insertStmt[len(lines)] = stmt
	}

	valueArgs := make([]interface{}, 0, len(lines)*nbFieldInsterted)

	for _, line := range lines {
		valueArgs = append(valueArgs, line.RemoteHost)
		valueArgs = append(valueArgs, line.Time)
		valueArgs = append(valueArgs, line.Request)
		valueArgs = append(valueArgs, line.Status)
		valueArgs = append(valueArgs, line.Bytes)
		valueArgs = append(valueArgs, line.Referer)
		valueArgs = append(valueArgs, line.UserAgent)
		valueArgs = append(valueArgs, line.URL)
	}
	_, err = stmt.ExecContext(ctx, valueArgs...)

	return err
}

func buildValuesQuery(nbFieldInsterted int, nbItem int) string {
	valuesQuery := make([]string, 0, nbItem)
	for i := 0; i < nbItem; i++ {
		oneInsert := make([]string, 0, nbFieldInsterted)
		for j := 0; j < nbFieldInsterted; j++ {
			oneInsert = append(oneInsert, fmt.Sprintf("$%d", i*nbFieldInsterted+j))
		}
		valuesQuery = append(valuesQuery, "("+strings.Join(oneInsert, ", ")+")")
	}
	return strings.Join(valuesQuery, ", ")
}

func (q *querier) Close() error {
	for _, stmt := range q.insertStmt {
		err := stmt.Close()
		if err != nil {
			return err
		}
	}
	return q.db.Close()
}
