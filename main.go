package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

func main() {
	searchEverywhere := flag.Bool("e", false, "Search everywhere, not just titles")
	captainsLog := flag.String("c", "", "Captain's Log: print ID of today's note, else <date>,<tag>")

	flag.Parse()

	db := NewDB()
	defer db.Close()

	if *captainsLog != "" {
		fmt.Println(CaptainsLog(db, *captainsLog))
	} else if flag.NArg() != 1 {
		log.Fatal("missing search term")
	} else {
		fmt.Println(Search(db, !*searchEverywhere, flag.Arg(0)))
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
