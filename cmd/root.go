package cmd

import (
	"github.com/mnadel/bearfred/cmd/journal"
	"github.com/mnadel/bearfred/cmd/search"
	"github.com/mnadel/bearfred/cmd/version"
	"github.com/spf13/cobra"
)

func Execute() error {
	rootCmd := &cobra.Command{
		Use:   "bearfred",
		Short: "A CLI for an Alfred+Bear integration",
		Long:  "Search notes, plus helpers to implement a daily journal",
	}

	rootCmd.AddCommand(journal.New())
	rootCmd.AddCommand(search.New())
	rootCmd.AddCommand(version.New())

	return rootCmd.Execute()
}
