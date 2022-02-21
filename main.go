package main

import (
	"log"

	"github.com/mnadel/bearfred/cmd/journal"
	"github.com/mnadel/bearfred/cmd/search"
	"github.com/mnadel/bearfred/cmd/version"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "bearfred",
		Short: "A CLI for an Alfred+Bear integration",
		Long:  "Search notes, plus helpers to implement a daily journal",
	}

	cmd.AddCommand(journal.New())
	cmd.AddCommand(search.New())
	cmd.AddCommand(version.New())

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
