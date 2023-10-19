package transcript

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mnadel/freddiebear/db"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	optAll      bool
	optShowTags bool
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transcript [tag]",
		Short: "Create a transcript for a tag",
		Long:  "Generate a date-based transcript for a given tag",
		Args:  cobra.ExactArgs(1),
		RunE:  runner,
	}

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
		transcript.WriteString(fmt.Sprintf("`%s`\n", result.ModificationDate))

		for _, note := range extractTaggedNote(result, args[0]) {
			transcript.WriteString(note)
			transcript.WriteString("\n")
		}
	}

	fmt.Printf(transcript.String())

	return nil
}

func extractTaggedNote(note *db.Record, tag string) []string {
	re := regexp.MustCompile("(?s)#"+tag+"\n(.*?)(?:##\\s(.*?)\\n|$)") // match after #tag and everything until the next ## heading
	matches := re.FindAllStringSubmatch(note.Text, -1)

	parts := make([]string, len(matches))

	for i := range matches {
		parts[i] = matches[i][1]
	}

	return parts
}
