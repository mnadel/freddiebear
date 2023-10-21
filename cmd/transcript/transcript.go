package transcript

import (
	"fmt"
	"strings"

	"github.com/mnadel/freddiebear/db"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	optAst   bool
	optDebug bool
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transcript [tag]",
		Short: "Create a transcript for a tag",
		Long:  "Generate a date-based transcript for a given tag",
		Args:  cobra.ExactArgs(1),
		RunE:  runner,
	}

	cmd.Flags().BoolVar(&optAst, "ast", false, "show representation of the AST")
	cmd.Flags().BoolVar(&optDebug, "debug", false, "print parsing debug info")

	return cmd
}

func runner(cmd *cobra.Command, args []string) error {
	bearDB, err := db.NewDB()
	if err != nil {
		return errors.WithStack(err)
	}
	defer bearDB.Close()

	results, err := bearDB.QueryTag(args[0])

	if err != nil {
		return errors.WithStack(err)
	}

	transcript := strings.Builder{}

	for _, result := range results {
		transcript.WriteString(fmt.Sprintf("## %s\n", result.Title))
		transcript.WriteString(fmt.Sprintf("_%s_\n", result.ModificationDate))

		te := NewTagExtractor([]byte(result.Text), "#"+args[0])
		data := te.ExtractTaggedNotes()

		transcript.Write(data)
	}

	if !optAst {
		fmt.Println(transcript.String())
	}

	return nil
}
