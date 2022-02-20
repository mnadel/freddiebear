package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"
)

func main() {
	searchEverywhere := flag.Bool("e", false, "Search everywhere, not just titles")
	captainsLog := flag.String("c", "", "Captain's Log: print ID of today's note, else <date>,<tag>")

	flag.Parse()

	db := NewDB()
	defer db.Close()

	if *captainsLog != "" {
		fmt.Println(doCaptainsLog(db, *captainsLog))
	} else if flag.NArg() != 1 {
		log.Fatal("missing search term")
	} else {
		fmt.Println(doSearch(db, !*searchEverywhere, flag.Arg(0)))
	}
}

func doSearch(db *DB, searchTitleOnly bool, searchTerm string) string {
	var results []Result
	var err error

	if searchTitleOnly {
		results, err = db.QueryTitles(searchTerm)
	} else {
		results, err = db.QueryText(searchTerm)
	}

	if err != nil {
		log.Fatal(err.Error())
	}

	return serialize(results)
}

func doCaptainsLog(db *DB, dateTag string) string {
	id, err := db.QueryCaptainsLog()
	if err != nil {
		log.Fatal(err.Error())
	}

	if id == "" {
		now := time.Now()
		return fmt.Sprintf("%s,%s/%s/%s", now.Format("2006-01-02"), dateTag, now.Format("2006"), now.Format("01"))
	} else {
		return id
	}
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
