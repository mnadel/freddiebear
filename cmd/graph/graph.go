package graph

import (
	"fmt"

	"github.com/mnadel/freddiebear/db"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	graphCmd := &cobra.Command{
		Use:   "graph",
		Short: "Visualize links between notes",
		Long:  "Generate a DOT graph of the links between notes",
		Args:  cobra.NoArgs,
		RunE:  runner,
	}

	return graphCmd
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

	fmt.Println("digraph Notes {")

	nodes := make(map[string]int)

	for i, edge := range graph {
		if _, ok := nodes[edge.SourceTitle]; !ok {
			nodes[edge.SourceTitle] = i
			fmt.Printf("	node_%d [label=\"%s\"];\n", i, edge.SourceTitle)
		}
		if _, ok := nodes[edge.DestinationTitle]; !ok {
			nodes[edge.DestinationTitle] = i + len(graph)
			fmt.Printf("	node_%d [label=\"%s\"];\n", i+len(graph), edge.DestinationTitle)
		}
	}

	fmt.Println("")

	for _, edge := range graph {
		src := nodes[edge.SourceTitle]
		dest := nodes[edge.DestinationTitle]
		fmt.Printf("	node_%d -> node_%d;\n", src, dest)
	}

	fmt.Println("}")

	return nil
}
