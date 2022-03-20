package export

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	"github.com/mnadel/freddiebear/db"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	searchCmd := &cobra.Command{
		Use:   "export [destination]",
		Short: "Export notes",
		Long:  "Export notes to Markdown files",
		Args:  cobra.ExactArgs(1),
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

	exporter, err := exporter(args[0])
	if err != nil {
		return errors.WithStack(err)
	}

	return bearDB.Export(exporter)
}

func exporter(destination string) (db.Exporter, error) {
	info, err := os.Stat(destination)
	if os.IsNotExist(err) {
		return nil, errors.WithStack(err)
	} else if !info.IsDir() {
		return nil, errors.WithStack(fmt.Errorf("not a directory: %s", destination))
	}

	return func(record *db.Record) error {
		filename := fmt.Sprintf("%s (%d).md", url.QueryEscape(record.Title), record.ID)
		filename = path.Join(destination, filename)

		if err = ioutil.WriteFile(filename, []byte(record.Text), 0644); err != nil {
			return err
		}

		return nil
	}, nil
}
