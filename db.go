package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"database/sql"
)

const (
	dbFile = "/Library/Group Containers/9K33E3U3T4.net.shinyfrog.bear/Application Data/database.sqlite"

	sqlCaptainsLog = `
		SELECT
			ZUNIQUEIDENTIFIER, ZTITLE
 		FROM
			ZSFNOTE
 		WHERE
			ZARCHIVED = 0
			AND ZTRASHED = 0
			AND ZTITLE = ?
	`

	sqlTitle = `
		SELECT DISTINCT
			ZUNIQUEIDENTIFIER, ZTITLE
 		FROM
			ZSFNOTE
 		WHERE
			ZARCHIVED = 0
			AND ZTRASHED = 0
			AND lower(ZTITLE) LIKE lower(?)
 		ORDER BY
			ZMODIFICATIONDATE DESC
	`

	sqlText = `
		SELECT DISTINCT
			ZUNIQUEIDENTIFIER, ZTITLE
 		FROM
			ZSFNOTE
 		WHERE
			ZARCHIVED = 0
			AND ZTRASHED = 0
			AND (lower(ZTEXT) LIKE lower(?) OR lower(ZTITLE) LIKE lower(?))
 		ORDER BY
			ZMODIFICATIONDATE DESC
	`
)

type DB struct {
	db *sql.DB
}

type Result struct {
	ID    string
	Title string
}

type Results []Result

func NewDB() *DB {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	dbFile := path.Join(home, dbFile)

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}

	pragmasSQL := `
		PRAGMA synchronous = normal;
		PRAGMA temp_store = memory;
		PRAGMA mmap_size = 30000000000;
		PRAGMA cache_size = -64000;
	`

	if _, err := db.Exec(pragmasSQL); err != nil {
		log.Fatal(err.Error())
	}

	return &DB{db}
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) QueryCaptainsLog() (string, error) {
	bind := time.Now().Format("2006-01-02")
	rows, err := d.db.Query(sqlCaptainsLog, bind)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	results, err := collectRows(rows)
	if err != nil {
		return "", err
	} else if len(results) > 1 {
		return "", fmt.Errorf("found too many records")
	} else if len(results) == 1 {
		return results[0].ID, nil
	} else {
		return "", nil
	}
}

func (d *DB) QueryTitles(term string) ([]Result, error) {
	bind := "%" + term + "%"
	rows, err := d.db.Query(sqlTitle, bind)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return collectRows(rows)
}

func (d *DB) QueryText(term string) ([]Result, error) {
	bind := "%" + term + "%"
	rows, err := d.db.Query(sqlText, bind, bind)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return collectRows(rows)
}

func collectRows(rows *sql.Rows) ([]Result, error) {
	var id string
	var title string

	results := make(Results, 0)

	for rows.Next() {
		err := rows.Scan(&id, &title)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, Result{ID: id, Title: title})
	}

	return results, rows.Err()
}
