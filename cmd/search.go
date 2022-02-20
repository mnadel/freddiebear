package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/mnadel/bearfred/db"
	"github.com/spf13/cobra"
)

var (
	searchAll bool

	searchCmd = &cobra.Command{
		Use:   "search <term>",
		Short: "Search for a note",
		Long:  "Generate search results in Alfred Workflow's XML schema format",
		Args:  cobra.ExactArgs(1),
		Run:   searchCmdRunner,
	}
)

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().BoolVar(&searchAll, "all", false, "search everything, else titles only")
}

func searchCmdRunner(cmd *cobra.Command, args []string) {
	bearDB := db.NewDB()
	defer bearDB.Close()

	var results []db.Result
	var err error

	if searchAll {
		results, err = bearDB.QueryText(args[0])
	} else {
		results, err = bearDB.QueryTitles(args[0], false)
	}

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Print(serialize(results))
}

func serialize(results db.Results) string {
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
