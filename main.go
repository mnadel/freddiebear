package main

import (
	"log"

	"github.com/mnadel/freddiebear/cmd/export"
	"github.com/mnadel/freddiebear/cmd/graph"
	"github.com/mnadel/freddiebear/cmd/journal"
	"github.com/mnadel/freddiebear/cmd/search"
	"github.com/mnadel/freddiebear/cmd/version"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "freddiebear",
		Short: "A CLI for an Alfred+Bear integration",
		Long:  "Search notes, plus helpers to implement a daily journal",
	}

	cmd.AddCommand(journal.New())
	cmd.AddCommand(search.New())
	cmd.AddCommand(version.New())
	cmd.AddCommand(export.New())
	cmd.AddCommand(graph.New())

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
