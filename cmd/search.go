package cmd

import (
	"fmt"
	"strings"

	"github.com/mnadel/bearfred/db"
	"github.com/spf13/cobra"
)

var (
	optSearchAll bool
)

func init() {
	searchCmd := &cobra.Command{
		Use:   "search [term]",
		Short: "Search for a note",
		Long:  "Generate search results in Alfred Workflow's XML schema format",
		Args:  cobra.ExactArgs(1),
		RunE:  searchCmdRunner,
	}

	searchCmd.Flags().BoolVar(&optSearchAll, "all", false, "search everything, else titles only")

	rootCmd.AddCommand(searchCmd)
}

func searchCmdRunner(cmd *cobra.Command, args []string) error {
	bearDB, err := db.NewDB()
	if err != nil {
		return err
	}
	defer bearDB.Close()

	var results db.Results

	if optSearchAll {
		results, err = bearDB.QueryText(args[0])
	} else {
		results, err = bearDB.QueryTitles(args[0], false)
	}

	if err != nil {
		return err
	}

	fmt.Print(serialize(results))

	return nil
}

func serialize(results db.Results) string {
	builder := strings.Builder{}

	builder.WriteString(`<?xml version="1.0" encoding="utf-8"?>`)
	builder.WriteString(`<items>`)

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
