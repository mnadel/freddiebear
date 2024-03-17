package forwardlinks

import (
	"fmt"
	"strings"

	"github.com/mnadel/freddiebear/alfred"
	"github.com/mnadel/freddiebear/db"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	searchCmd := &cobra.Command{
		Use:   "forwardlinks [term]",
		Short: "Show forward links for notes matching search term",
		Long:  "Generate forward link results in Alfred Workflow's XML schema format",
		Args:  cobra.ExactArgs(1),
		RunE:  runner,
	}

	return searchCmd
}

func runner(cmd *cobra.Command, args []string) error {
	bearDB, err := db.NewDB()
	if err != nil {
		return errors.WithStack(err)
	}
	defer bearDB.Close()

	graph, err := bearDB.QueryGraph()
	if err != nil {
		return errors.WithStack(err)
	}

	term := strings.ToLower(args[0])
	matches := make(map[*db.Result]*db.Result)

	for _, edge := range graph {
		if strings.Contains(strings.ToLower(edge.Source.Title), term) {
			matches[edge.Source] = edge.Target
		}
	}

	fmt.Println(alfred.AlfredBacklinkXML(matches))

	return nil
}
