package backlinks

import (
	"fmt"
	"strings"

	"github.com/mnadel/freddiebear/cmd/alfred"
	"github.com/mnadel/freddiebear/db"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	searchCmd := &cobra.Command{
		Use:   "backlinks [term]",
		Short: "Show backlinks for notes matching search term",
		Long:  "Generate backlink results in Alfred Workflow's XML schema format",
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
		if strings.Contains(strings.ToLower(edge.Target.Title), term) {
			matches[edge.Target] = edge.Source
		}
	}

	fmt.Println(alfred.AlfredBacklinkXML(matches))

	return nil
}
