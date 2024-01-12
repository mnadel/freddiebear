package cleanup

import (
	"fmt"

	"github.com/mnadel/freddiebear/db"
	"github.com/mnadel/freddiebear/cmd/export"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	searchCmd := &cobra.Command{
		Use:   "cleanup",
		Short: "Cleanup your export",
		Long:  "Generate list of attachments that can be deleted",
		Args:  cobra.ExactArgs(0),
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

	attachments, err := bearDB.QueryDeletedAttachments()
	if err != nil {
		return errors.WithStack(err)
	}

	for _, attach := range attachments {
		fmt.Println(export.BuildAttachmentFilename(attach))
	}

	return nil
}
