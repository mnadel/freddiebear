package db

import (
	"os"
	"path"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"

	"database/sql"
)

const (
	dbFile = `/Library/Group Containers/9K33E3U3T4.net.shinyfrog.bear/Application Data/database.sqlite`

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

	sqlPragma = `
		PRAGMA query_only = on;
		PRAGMA synchronous = normal;
		PRAGMA temp_store = memory;
		PRAGMA mmap_size = 30000000000;
		PRAGMA cache_size = -64000;
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

// Results is a list of *Result, and represents a collection of notes in the database
type Results []*Result

// Create a new DB, referencing the user's Bear Notes database
func NewDB() (*DB, error) {
	url, err := dbUrl()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	db, err := sql.Open("sqlite3", url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if _, err := db.Exec(sqlPragma); err != nil {
		return nil, errors.WithStack(err)
	}

	return &DB{db}, nil
}

// Close cleans up our database connection
func (d *DB) Close() error {
	return d.db.Close()
}

// QueryTitles searches for a term within the titles of notes within the database, setting
// `exact` to true will do an exact match, else it'll perform a substring match
func (d *DB) QueryTitles(term string, exact bool) (Results, error) {
	var bind string

	if exact {
		bind = term
	} else {
		bind = substringSearch(term)
	}

	rows, err := d.db.Query(sqlTitle, bind)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	return rowsToResults(rows)
}

// QueryText searches for a term within the body or title of notes within the database
func (d *DB) QueryText(term string) (Results, error) {
	bind := substringSearch(term)
	rows, err := d.db.Query(sqlText, bind, bind)
	if err != nil {
		return nil, errors.WithStack(err)
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
			return nil, errors.WithStack(err)
		}
		results = append(results, &Result{ID: id, Title: title})
	}

	return results, errors.WithStack(rows.Err())
}

func substringSearch(term string) string {
	bind := strings.Builder{}
	bind.WriteString(`%`)
	bind.WriteString(term)
	bind.WriteString(`%`)
	return bind.String()
}

func dbUrl() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.WithStack(err)
	}

	url := strings.Builder{}
	url.WriteString(path.Join(home, dbFile))
	url.WriteString(`?mode=ro`)

	return url.String(), nil
}
