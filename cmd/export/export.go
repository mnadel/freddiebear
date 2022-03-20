package export

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	"github.com/mnadel/freddiebear/db"
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
		return bearDB.Export(func(rec *db.Record) error {
			fmt.Println(path.Join(args[0], buildFilename(rec)))
			return nil
		})
	}

	info, err := os.Stat(args[0])
	if os.IsNotExist(err) {
		return errors.WithStack(err)
	} else if !info.IsDir() {
		return errors.WithStack(fmt.Errorf("not a directory: %s", args[0]))
	}

	exporter, err := exporter(args[0])
	if err != nil {
		return errors.WithStack(err)
	}

	return bearDB.Export(exporter)
}

func exporter(destinationDir string) (db.Exporter, error) {
	return func(record *db.Record) error {
		filename := buildFilename(record)
		outfile := path.Join(destinationDir, filename)

		if err := ioutil.WriteFile(outfile, []byte(record.Text), 0644); err != nil {
			return err
		}

		return nil
	}, nil
}

func buildFilename(record *db.Record) string {
	id := fmt.Sprintf("%x", md5.Sum([]byte(record.GUID)))
	filename := fmt.Sprintf("%s (%s).md", url.QueryEscape(record.Title), id[0:7])

	return filename
}
