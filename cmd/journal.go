package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/mnadel/bearfred/db"
	"github.com/spf13/cobra"
)

var (
	tagName       string
	tagAppendDate bool

	journalCmd = &cobra.Command{
		Use:   "journal",
		Short: "Daily journal helper",
		Long:  "Display daily note ID, or <title>,<tag>",
		Run:   journalCmdRunner,
	}
)

func init() {
	rootCmd.AddCommand(journalCmd)
	journalCmd.Flags().StringVar(&tagName, "tag", "", "tag to add to journal entry")
	journalCmd.Flags().BoolVar(&tagAppendDate, "date", false, "append date (/yyyy/mm) to tag")
}

func journalCmdRunner(cmd *cobra.Command, args []string) {
	bearDB := db.NewDB()
	defer bearDB.Close()

	now := time.Now()
	term := now.Format("2006-01-02")

	results, err := bearDB.QueryTitles(term, true)
	if err != nil {
		log.Fatal(err.Error())
	}

	var id string

	if len(results) > 1 {
		log.Fatal("found too many matches")
	} else if len(results) == 1 {
		id = results[0].ID
	}

	if id == "" {
		tag := journalTag(now)
		fmt.Printf("%s,%s", now.Format("2006-01-02"), tag)
	} else {
		fmt.Print(id)
	}
}

func journalTag(now time.Time) string {
	if tagName == "" && !tagAppendDate {
		return ""
	} else if !tagAppendDate {
		return tagName
	} else {
		return fmt.Sprintf("%s/%s/%s", tagName, now.Format("2006"), now.Format("01"))
	}
}