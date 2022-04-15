package search

import (
	"fmt"
	"strings"

	"github.com/mnadel/freddiebear/db"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	optAll      bool
	optShowTags bool
)

func New() *cobra.Command {
	searchCmd := &cobra.Command{
		Use:   "search [term]",
		Short: "Search for a note",
		Long:  "Generate search results in Alfred Workflow's XML schema format",
		Args:  cobra.ExactArgs(1),
		RunE:  runner,
	}

	searchCmd.Flags().BoolVar(&optAll, "all", false, "search everything (default: titles only)")
	searchCmd.Flags().BoolVar(&optShowTags, "show-tags", false, "include tags in output")

	return searchCmd
}

func runner(cmd *cobra.Command, args []string) error {
	bearDB, err := db.NewDB()
	if err != nil {
		return errors.WithStack(err)
	}
	defer bearDB.Close()

	var results db.Results

	if optAll {
		results, err = bearDB.QueryText(args[0])
	} else {
		results, err = bearDB.QueryTitles(args[0], false)
	}

	if err != nil {
		return errors.WithStack(err)
	}

	fmt.Print(serialize(results, args[0]))

	return nil
}

func serialize(results db.Results, searchTerm string) string {
	builder := strings.Builder{}

	builder.WriteString(`<?xml version="1.0" encoding="utf-8"?>`)
	builder.WriteString(`<items>`)

	if len(results) == 0 {
		builder.WriteString(`<item valid="yes">`)
		builder.WriteString(`<subtitle>Create note</subtitle>`)
		builder.WriteString(`<title>`)
		builder.WriteString(searchTerm)
		builder.WriteString(`</title>`)
		builder.WriteString(`<arg>`)
		builder.WriteString(`create:`)
		builder.WriteString(searchTerm)
		builder.WriteString(`</arg>`)
		builder.WriteString(`</item>`)
	} else {
		for _, item := range results {
			builder.WriteString(`<item valid="yes">`)
			builder.WriteString(`<title>`)
			builder.WriteString(item.Title)
			builder.WriteString(`</title>`)

			if !optShowTags {
				builder.WriteString(`<subtitle>Open note</subtitle>`)
			} else {
				builder.WriteString(`<subtitle>`)
				builder.WriteString(strings.Join(item.UniqueTags(), ", "))
				builder.WriteString(`</subtitle>`)
			}

			builder.WriteString(`<arg>`)
			builder.WriteString(item.ID)
			builder.WriteString(`</arg>`)
			builder.WriteString(`</item>`)
		}
	}

	builder.WriteString(`</items>`)

	return builder.String()
}
