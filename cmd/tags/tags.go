package tags

import (
	"fmt"
	"strings"

	"github.com/mnadel/freddiebear/db"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags",
		Short: "Generate a list of all tags",
		Long:  "Generate a list of all tags in Alfred Workflow's JSON schema format",
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

	allTags, err := bearDB.QueryTags()
	if err != nil {
		return errors.WithStack(err)
	}

	items := make([]string, len(allTags))

	for i, t := range allTags {
		items[i] = fmt.Sprintf(`{"title":"%s","arg":"%s"}`, t, t)
	}

	fmt.Printf(`{"items":[%s]}`, strings.Join(items, ","))

	return nil
}
