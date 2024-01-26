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
		Use:   "graph [term]",
		Short: "Visualize links between notes",
		Long:  "Generate a DOT graph of the links between notes. If term specified source or target must contain it.",
		Args:  cobra.MaximumNArgs(1),
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
	includedEdges := make(map[*db.Edge]bool)

	for i, edge := range graph {
		if len(args) == 1 && args[0] != "" {
			haystack := strings.Builder{}
			haystack.WriteString(strings.ToLower(edge.Source.Title))
			haystack.WriteString(strings.ToLower(edge.Target.Title))

			if !strings.Contains(haystack.String(), strings.ToLower(args[0])) {
				continue
			}
		}

		includedEdges[edge] = true

		if _, ok := nodes[edge.Source.Title]; !ok {
			nodes[edge.Source.Title] = i

			fmt.Printf("	node_%d [label=%s];\n", i, nodeLabel(edge.Source))
		}
		if _, ok := nodes[edge.Target.Title]; !ok {
			uniqID := i + len(graph)
			nodes[edge.Target.Title] = uniqID

			fmt.Printf("	node_%d [label=%s];\n", uniqID, nodeLabel(edge.Target))
		}
	}

	fmt.Println("")

	for edge, _ := range includedEdges {
		src := nodes[edge.Source.Title]
		dest := nodes[edge.Target.Title]
		fmt.Printf("	node_%d -> node_%d;\n", dest, src)
	}

	fmt.Println("}")

	return nil
}

func nodeLabel(n *db.Result) string {
	tags := n.UniqueTags()
	if len(tags) > 0 {
		tags[0] = fmt.Sprintf("#%s", tags[0])
	} else {
		tags = []string{"&nbsp;"}
	}

	alltags := strings.Join(tags, ", #")

	return fmt.Sprintf(`<
		%s
		<br/>
		<font point-size="10">%s</font>
	>`, util.ToSafeString(n.Title), alltags)
}
