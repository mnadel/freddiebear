package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

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

	items := make([][]string, 0)

	for rows.Next() {
		err := rows.Scan(&id, &title)
		if err != nil {
			log.Fatal(err)
		}
		items = append(items, []string{title, id})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	serialize(items)
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

func serialize(results [][]string) {
	fmt.Println(`<?xml version="1.0" encoding="utf-8"?><items>`)

	for _, item := range results {
		fmt.Println(`<item valid="yes">` +
			`<title>` + item[0] + `</title><subtitle>Open note</subtitle>` +
			`<arg>` + item[1] + `</arg></item>`)
	}

	fmt.Println(`</items>`)
}
