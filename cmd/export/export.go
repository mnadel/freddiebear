package export

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/mnadel/freddiebear/db"
	"github.com/mnadel/freddiebear/db/exporter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	listOnly bool
)

func New() *cobra.Command {
	searchCmd := &cobra.Command{
		Use:   "export [destination]",
		Short: "Export notes",
		Long:  "Export notes to Markdown files",
		Args:  cobra.ExactArgs(1),
		RunE:  runner,
	}

	searchCmd.Flags().BoolVar(&listOnly, "list", false, "list files to export, but don't create them")

	return searchCmd
}

func runner(cmd *cobra.Command, args []string) error {
	bearDB, err := db.NewDB()
	if err != nil {
		return errors.WithStack(err)
	}
	defer bearDB.Close()

	if listOnly {
		return bearDB.Export(printingExporter(args[0]))
	}

	info, err := os.Stat(args[0])
	if os.IsNotExist(err) {
		return errors.WithStack(err)
	} else if !info.IsDir() {
		return errors.WithStack(fmt.Errorf("not a directory: %s", args[0]))
	}

	exp, err := writingExporter(args[0])
	if err != nil {
		return errors.WithStack(err)
	}

	return bearDB.Export(exp)
}

func printingExporter(destinationDir string) db.Exporter {
	return func(record *db.Record) error {
		fmt.Println(path.Join(destinationDir, exporter.BuildFilename(record)))
		return nil
	}
}

func writingExporter(destinationDir string) (db.Exporter, error) {
	exp, err := exporter.NewExporter(destinationDir)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return func(record *db.Record) error {
		if renamed, oldName := exp.IsRenamed(record); renamed {
			log.Println("detected rename of", oldName)

			if err := os.Remove(string(oldName)); err != nil {
				return errors.WithStack(err)
			}
		} else {
			changes, err := exp.IsChanged(record)
			if err != nil {
				return errors.WithStack(err)
			} else if !changes {
				return nil
			}
		}

		log.Println("exporting", record.SHA, record.Title)
		return writeRecord(record, destinationDir)
	}, nil
}

func writeRecord(record *db.Record, destinationDir string) error {
	filename := path.Join(destinationDir, exporter.BuildFilename(record))

	if err := os.WriteFile(filename, []byte(record.Text), 0644); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
