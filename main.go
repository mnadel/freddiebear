package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DBFILE = "/Library/Group Containers/9K33E3U3T4.net.shinyfrog.bear/Application Data/database.sqlite"

	SQL_TITLE = `
		SELECT DISTINCT
			ZUNIQUEIDENTIFIER, ZTITLE
 		FROM
			ZSFNOTE
 		WHERE
			ZARCHIVED=0
			AND ZTRASHED=0
			AND lower(ZTITLE) LIKE lower(?)
 		ORDER BY
			ZMODIFICATIONDATE DESC
	`

	SQL_TEXT = `
		SELECT DISTINCT
			ZUNIQUEIDENTIFIER, ZTITLE
 		FROM
			ZSFNOTE
 		WHERE
			ZARCHIVED=0
			AND ZTRASHED=0
			AND lower(ZTEXT) LIKE lower(?)
 		ORDER BY
			ZMODIFICATIONDATE DESC
	`
)

type Result struct {
	ID    string
	Title string
}

type Results []Result

var (
	searchEverywhere bool
	searchTerm       string
)

func init() {
	flag.BoolVar(&searchEverywhere, "e", false, "Search everywhere, not just titles")

	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal("missing search term")
	} else {
		searchTerm = flag.Arg(0)
	}
}

func main() {
	db := openDB()
	defer db.Close()

	sql := SQL_TITLE
	if searchEverywhere {
		sql = SQL_TEXT
	}

	rows, err := db.Query(sql, "%"+searchTerm+"%")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

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
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(serialize(results))
}

func openDB() *sql.DB {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	dbFile := path.Join(home, DBFILE)

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func serialize(results Results) string {
	builder := strings.Builder{}

	builder.WriteString(`<?xml version="1.0" encoding="utf-8"?><items>`)

	for _, item := range results {
		builder.WriteString(`<item valid="yes">`)
		builder.WriteString(`<title>`)
		builder.WriteString(item.Title)
		builder.WriteString(`</title>`)
		builder.WriteString(`<subtitle>Open note</subtitle>`)
		builder.WriteString(`<arg>`)
		builder.WriteString(item.ID)
		builder.WriteString(`</arg>`)
		builder.WriteString(`</item>`)
	}

	builder.WriteString(`</items>`)

	return builder.String()
}
