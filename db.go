package main

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var wg sync.WaitGroup

type dbInsertResult struct {
	log clf
	err error
}

func dbConnect() (db *sql.DB, err error) {
	dbHost := os.Getenv("DBHOST")
	dbPort := os.Getenv("DBPORT")
	dbUser := os.Getenv("DBUSER")
	dbPass := os.Getenv("DBPASS")
	dbDatabase := os.Getenv("DBDATABASE")

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbDatabase,
	)
	db, err = sql.Open("postgres", connStr)

	return
}

func dbInsertLog(c chan dbInsertResult, l clf, db *sql.DB) {
	defer wg.Done()

	res := new(dbInsertResult)
	res.log = l

	q := `INSERT INTO
		logs (
			client_host,
			rfc1413,
			remote_user,
			date_time,
			method,
			resource,
			protocol,
			status,
			size,
			referer,
			user_agent
		)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT
			DO NOTHING`

	_, res.err = db.Exec(q,
		l.ClientHost,
		l.RFC1413,
		l.RemoteUser,
		l.DateTime.Time,
		l.Method,
		l.Resource,
		l.Protocol,
		l.Status,
		l.Size,
		l.Referer,
		l.UserAgent,
	)

	c <- *res
	return
}

func dbSelectLastLogs(method string, n int64, db *sql.DB) (logs []clf, err error) {
	q := `SELECT
			client_host,
			rfc1413,
			remote_user,
			date_time,
			method,
			resource,
			protocol,
			status,
			size,
			referer,
			user_agent
		FROM
			(SELECT
				id,
				client_host,
				rfc1413,
				remote_user,
				date_time,
				method,
				resource,
				protocol,
				status,
				size,
				referer,
				user_agent
			FROM
				logs
			WHERE
				method = $1
			ORDER BY
				id
			DESC
			LIMIT $2) alias
		ORDER BY id ASC`
	// It would be necessary to cast the result of the subtraction (OFFSET) to unsigned.
	// Unfortunately postgresql doesn't have unsigned types. Then ORDER BY was used.
	// OFFSET
	// (SELECT COUNT(*) FROM logs WHERE method = $1) - $2

	rows, err := db.Query(q, method, n)
	if err != nil {
		return []clf{}, err
	}
	defer rows.Close()

	var l clf
	for rows.Next() {
		err = rows.Scan(
			&l.ClientHost,
			&l.RFC1413,
			&l.RemoteUser,
			&l.DateTime.Time,
			&l.Method,
			&l.Resource,
			&l.Protocol,
			&l.Status,
			&l.Size,
			&l.Referer,
			&l.UserAgent,
		)
		if err != nil {
			return
		}
		logs = append(logs, l)
	}
	err = rows.Err()
	return
}
