package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "bearfred",
	Short: "A CLI for an Alfred+Bear integration",
	Long:  `Search notes and helpers to implement a daily journal`,
}

func Execute() error {
	return rootCmd.Execute()
}
