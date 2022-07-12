package db

import (
	"crypto/md5"
	"fmt"
	"os"
	"path"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mnadel/freddiebear/util"
	"github.com/pkg/errors"

	"database/sql"
)

const (
	dbFile = `/Library/Group Containers/9K33E3U3T4.net.shinyfrog.bear/Application Data/database.sqlite?mode=ro`

	sqlTitle = `
		SELECT DISTINCT
			note.ZUNIQUEIDENTIFIER,
			note.ZTITLE,
			GROUP_CONCAT(COALESCE(tag.ZTITLE, ''))
		FROM
			ZSFNOTE note
			LEFT OUTER JOIN Z_7TAGS tags ON note.Z_PK = tags.Z_7NOTES
			LEFT OUTER JOIN ZSFNOTETAG tag ON tags.Z_14TAGS = tag.Z_PK
		WHERE
			note.ZARCHIVED = 0
			AND note.ZTRASHED = 0
			AND LOWER(note.ZTITLE) LIKE LOWER(?)
		GROUP BY
			note.ZUNIQUEIDENTIFIER
		ORDER BY
			note.ZMODIFICATIONDATE DESC
	`

	sqlText = `
		SELECT DISTINCT
			note.ZUNIQUEIDENTIFIER,
			note.ZTITLE,
			GROUP_CONCAT(COALESCE(tag.ZTITLE, ''))
		FROM
			ZSFNOTE note
			LEFT OUTER JOIN Z_7TAGS tags ON note.Z_PK = tags.Z_7NOTES
			LEFT OUTER JOIN ZSFNOTETAG tag ON tags.Z_14TAGS = tag.Z_PK
		WHERE
			note.ZARCHIVED = 0
			AND note.ZTRASHED = 0
			AND (LOWER(note.ZTEXT) LIKE LOWER(?) OR LOWER(note.ZTITLE) LIKE LOWER(?))
		GROUP BY
			note.ZUNIQUEIDENTIFIER
		ORDER BY
			note.ZMODIFICATIONDATE DESC
	`

	sqlExport = `
		select
			ZUNIQUEIDENTIFIER,
			ZTITLE,
			ZTEXT
		from
			ZSFNOTE
		where
			ZARCHIVED = 0 
			and ZTRASHED = 0
	`

	sqlPragma = `
		PRAGMA query_only = on;
		PRAGMA synchronous = off;
		PRAGMA mmap_size = 30000000000;
		PRAGMA temp_store = memory;
		PRAGMA journal_mode = off;
		PRAGMA page_size = 512;
		PRAGMA locking_mode = exclusive;
	`
)

// Exporter is a func that receives an exported record
type Exporter func(record *Record) error

// DB represents the Bear Notes database
type DB struct {
	db *sql.DB
}

// Record represents an exported note
type Record struct {
	SHA   string
	Title string
	Text  string
}

// Result references a specific note: its identifier and title
type Result struct {
	ID    string
	Title string
	Tags  string
}

// Results is a list of *Result, and represents a collection of notes in the database
type Results []*Result

// Create a new DB, referencing the user's Bear Notes database
func NewDB() (*DB, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	db, err := sql.Open("sqlite3", path.Join(home, dbFile))
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

// Records returns the list of notes in the database
func (d *DB) Records() ([]*Record, error) {
	records := make([]*Record, 0)

	rows, err := d.db.Query(sqlExport)
	if err != nil {
		return nil, errors.WithStack(rows.Err())
	}

	var guid, title, text string

	for rows.Next() {
		err := rows.Scan(&guid, &title, &text)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		record := &Record{
			SHA:   fmt.Sprintf("%x", md5.Sum([]byte(guid)))[0:7],
			Title: title,
			Text:  text,
		}

		records = append(records, record)
	}

	return records, nil
}

// Export notes to specified directory
func (d *DB) Export(exporter Exporter) error {
	records, err := d.Records()
	if err != nil {
		return errors.WithStack(err)
	}

	for _, record := range records {
		if err = exporter(record); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
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

// UniqueTags returns the leaf-node tags ([a a/b a/b/c d] -> [a/b/c d])
func (r *Result) UniqueTags() []string {
	split := strings.Split(r.Tags, ",")
	return util.RemoveIntermediatePrefixes(split, "/")
}

// TitleCase returns a Alfred-safe version of the proper title casing
func (r *Result) TitleCase() string {
	return util.ToSafeString(util.ToTitleCase(r.Title))
}

func rowsToResults(rows *sql.Rows) (Results, error) {
	var id string
	var title string
	var tags string

	results := make(Results, 0)

	for rows.Next() {
		err := rows.Scan(&id, &title, &tags)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		results = append(results, &Result{ID: id, Title: title, Tags: tags})
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
