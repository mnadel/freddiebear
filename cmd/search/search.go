package search

import (
	"fmt"

	"github.com/mnadel/freddiebear/alfred"
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

	searchCmd.Flags().BoolVar(&optAll, "all", false, "full text search (default: titles only)")
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

	if len(results) == 0 {
		fmt.Print(alfred.AlfredCreateXML(args[0]))
	} else {
		fmt.Print(alfred.AlfredOpenXML(results, optShowTags))
	}

	return nil
}
