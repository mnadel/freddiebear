package db

import (
	"log"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"

	"database/sql"
)

const (
	dbFile = "/Library/Group Containers/9K33E3U3T4.net.shinyfrog.bear/Application Data/database.sqlite"

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

// DB represents the Bear Notes database
type DB struct {
	db *sql.DB
}

// Result references a specific note: its identifier and title
type Result struct {
	ID    string
	Title string
}

// Results is a collection of Result, and represents a set of notes in the database
type Results []Result

// Create a new DB, referencing the user's Bear Notes database
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

// Close cleans up our database connection
func (d *DB) Close() error {
	return d.db.Close()
}

// QueryTitles searches for a term within the titles of notes within the database
func (d *DB) QueryTitles(term string, exact bool) (Results, error) {
	bind := "%" + term + "%"
	if exact {
		bind = term
	}

	rows, err := d.db.Query(sqlTitle, bind)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rowsToResults(rows)
}

// QueryText searches for a term within the body or title of notes within the database
func (d *DB) QueryText(term string) (Results, error) {
	bind := "%" + term + "%"
	rows, err := d.db.Query(sqlText, bind, bind)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return rowsToResults(rows)
}

func rowsToResults(rows *sql.Rows) (Results, error) {
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
