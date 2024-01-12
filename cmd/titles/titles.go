package titles

import (
	"fmt"
	"strings"

	"github.com/mnadel/freddiebear/db"
	"github.com/mnadel/freddiebear/util"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "titles",
		Short: "Generate a list of all titles",
		Long:  "Generate a list of all titles in Alfred Workflow's JSON schema format",
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

	allTitles, err := bearDB.QueryAllTitles()
	if err != nil {
		return errors.WithStack(err)
	}

	items := make([]string, 0)

	for _, t := range allTitles {
		if !strings.Contains(t.Tags, "captainslog") {
			tags := strings.Split(t.Tags, ",")
			filtered := util.RemoveIntermediatePrefixes(tags, "/")
			tag := strings.Join(filtered, ",")
			items = append(items, fmt.Sprintf(`{"title":"%s","arg":"%s","subtitle":"%s"}`, t.Title, t.ID, tag))
		}
	}

	fmt.Printf(`{"items":[%s]}`, strings.Join(items, ","))

	return nil
}
