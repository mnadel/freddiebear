package graph

import (
	"fmt"
	"strings"

	"github.com/mnadel/freddiebear/db"
	"github.com/mnadel/freddiebear/util"
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
		if _, ok := nodes[edge.Source.Title]; !ok {
			nodes[edge.Source.Title] = i
			label := nodeLabel(edge.Source)
			fmt.Printf("	node_%d [label=%s];\n", i, label)
		}
		if _, ok := nodes[edge.Target.Title]; !ok {
			nodes[edge.Target.Title] = i + len(graph)
			label := nodeLabel(edge.Target)
			fmt.Printf("	node_%d [label=%s];\n", i+len(graph), label)
		}
	}

	fmt.Println("")

	for _, edge := range graph {
		src := nodes[edge.Source.Title]
		dest := nodes[edge.Target.Title]
		fmt.Printf("	node_%d -> node_%d;\n", src, dest)
	}

	fmt.Println("}")

	return nil
}

func nodeLabel(n *db.Node) string {
	alltags := strings.Join(n.UniqueTags(), ", ")
	return fmt.Sprintf(`<
		<table border="0" cellborder="0"><tr><td>%s</td></tr><tr><td><font size="-2">%s</font></td></tr></table>
	>`, util.ToSafeString(n.Title), alltags)
}
