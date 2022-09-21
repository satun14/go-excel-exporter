package main

import (
	"database/sql"
	"strings"
)

func execSqlQuery(request Request) (*sql.Rows, error) {
	db, err := sql.Open("mysql",
		strings.Join([]string{
			request.Db.User,
			":",
			request.Db.Password,
			"@tcp(",
			request.Db.Host,
			":",
			request.Db.Port,
			")/",
			request.Db.Name,
			"?parseTime=true"}, ""))

	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query(request.Sql)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
