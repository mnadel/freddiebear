package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	searchEverywhere bool
	searchTerm       string
	todaysNote       bool
)

func init() {
	flag.BoolVar(&searchEverywhere, "e", false, "Search everywhere, not just titles")
	flag.BoolVar(&todaysNote, "t", false, "Print ID of today's note if it exists")

	flag.Parse()
}

func main() {
	db := NewDB()
	defer db.Close()

	if todaysNote {
		id, err := db.QueryToday()
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(id)
		os.Exit(0)
	}

	if flag.NArg() != 1 {
		log.Fatal("missing search term")
	} else {
		searchTerm = flag.Arg(0)
	}

	var results []Result
	var err error

	if searchEverywhere {
		results, err = db.QueryText(searchTerm)
	} else {
		results, err = db.QueryTitles(searchTerm)
	}

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(serialize(results))
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
