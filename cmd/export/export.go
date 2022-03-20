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

var (
	optAll      bool
	optShowTags bool
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

	return func(title, text string) error {
		filename := fmt.Sprintf("%s.md", url.QueryEscape(title))
		filename = path.Join(destination, filename)

		if err = ioutil.WriteFile(filename, []byte(text), 0644); err != nil {
			return err
		}

		return nil
	}, nil
}
