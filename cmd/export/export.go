package export

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"regexp"

	"github.com/mnadel/freddiebear/db"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	listOnly bool
)

type SHA string
type Filename string

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

	exporter, err := writingExporter(args[0])
	if err != nil {
		return errors.WithStack(err)
	}

	return bearDB.Export(exporter)
}

func printingExporter(destinationDir string) db.Exporter {
	return func(record *db.Record) error {
		fmt.Println(path.Join(destinationDir, buildFilename(record)))
		return nil
	}
}

func writingExporter(destinationDir string) (db.Exporter, error) {
	filenameSHAs, err := getFilenameSHAs(destinationDir)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return func(record *db.Record) error {
		preferredFilename := path.Join(destinationDir, buildFilename(record))
		safeFilename := path.Join(destinationDir, buildSafeFilename(record))
		recordText := []byte(record.Text)

		changes, err := hasChanges(filenameSHAs, record.SHA, recordText)
		if err != nil {
			return errors.WithStack(err)
		} else if !changes {
			return nil
		}

		log.Println("exporing", record.SHA)

		if err := os.WriteFile(preferredFilename, recordText, 0644); err != nil {
			if err = os.WriteFile(safeFilename, recordText, 0644); err != nil {
				return errors.WithStack(err)
			}
		}

		return nil
	}, nil
}

func buildFilename(record *db.Record) string {
	return fmt.Sprintf("%s (%s).md", record.Title, record.SHA)
}

func buildSafeFilename(record *db.Record) string {
	return fmt.Sprintf("%s (%s).md", url.QueryEscape(record.Title), record.SHA)
}

func hasChanges(mapping map[SHA]Filename, sha string, newData []byte) (bool, error) {
	filename, ok := mapping[SHA(sha)]
	if !ok {
		return true, nil
	}

	oldData, err := os.ReadFile(string(filename))
	if err != nil {
		return false, errors.WithStack(err)
	}

	oldSum := md5.Sum(oldData)
	newSum := md5.Sum(newData)

	return oldSum != newSum, nil
}

func getFilenameSHAs(directory string) (map[SHA]Filename, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	filenames := make(map[SHA]Filename)
	re := regexp.MustCompile(`.*\s\((\w+)\)\.md$`)

	for _, file := range files {
		if !file.IsDir() {
			parts := re.FindStringSubmatch(file.Name())
			if len(parts) == 2 {
				filenames[SHA(parts[1])] = Filename(file.Name())
			}
		}
	}

	return filenames, nil
}
