package journal

import (
	"fmt"
	"time"

	"github.com/mnadel/bearfred/db"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	optTagName       string
	optTagAppendDate bool
)

func New() *cobra.Command {
	journalCmd := &cobra.Command{
		Use:   "journal",
		Short: "Daily journal helper",
		Long:  "Display daily note ID, or <title>,<tag>",
		RunE:  runner,
	}

	journalCmd.Flags().StringVar(&optTagName, "tag", "", "tag to add to journal entry")
	journalCmd.Flags().BoolVar(&optTagAppendDate, "date", false, "append date (yyyy/mm) to tag")

	return journalCmd
}

func runner(cmd *cobra.Command, args []string) error {
	bearDB, err := db.NewDB()
	if err != nil {
		return errors.WithStack(err)
	}
	defer bearDB.Close()

	now := time.Now()
	term := now.Format("2006-01-02")

	results, err := bearDB.QueryTitles(term, true)
	if err != nil {
		return errors.WithStack(err)
	}

	var id string

	if len(results) > 1 {
		return fmt.Errorf("found too many matches")
	} else if len(results) == 1 {
		id = results[0].ID
	}

	if id == "" {
		fmt.Printf("%s,%s", now.Format("2006-01-02"), journalTag(now))
	} else {
		fmt.Print(id)
	}

	return nil
}

func journalTag(now time.Time) string {
	if optTagName == "" && !optTagAppendDate {
		return ""
	} else if !optTagAppendDate {
		return optTagName
	} else {
		return fmt.Sprintf("%s/%s/%s", optTagName, now.Format("2006"), now.Format("01"))
	}
}
