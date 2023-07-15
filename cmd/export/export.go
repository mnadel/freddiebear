package export

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/mnadel/freddiebear/db"
	"github.com/mnadel/freddiebear/db/exporter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	RelativeTrashDirectoryPath = "Trash"
)

var (
	preview bool
	list    bool
)

func New() *cobra.Command {
	searchCmd := &cobra.Command{
		Use:   "export [destination]",
		Short: "Export notes",
		Long:  "Export notes to Markdown files",
		Args:  cobra.ExactArgs(1),
		RunE:  runner,
	}

	searchCmd.Flags().BoolVar(&preview, "preview", false, "list files that would be exported")
	searchCmd.Flags().BoolVar(&list, "list", false, "list files in export directory")

	return searchCmd
}

func runner(cmd *cobra.Command, args []string) error {
	bearDB, err := db.NewDB()
	if err != nil {
		return errors.WithStack(err)
	}
	defer bearDB.Close()

	if preview {
		return bearDB.Export(printingExporter(args[0]))
	} else if list {
		files, err := exporter.ListFiles(args[0])
		if err != nil {
			return errors.WithStack(err)
		}

		for _, f := range files {
			fmt.Println(f)
		}

		return nil
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

	if err := bearDB.Export(exp); err != nil {
		return errors.WithStack(err)
	}

	records, err := bearDB.Records()
	if err != nil {
		return errors.WithStack(err)
	}

	exporter, err := exporter.NewExporter(args[0])
	if err != nil {
		return errors.WithStack(err)
	}

	trashDir := path.Join(args[0], RelativeTrashDirectoryPath)
	_, err = os.Stat(trashDir)
	if os.IsNotExist(err) {
		if err := os.Mkdir(trashDir, 0755); err != nil {
			return errors.WithStack(err)
		}
	}

	if err := writeAttachmentMappings(args[0], bearDB); err != nil {
		errors.WithStack(err)
	}

	return exporter.Archive(records, trashDir)
}

func writeAttachmentMappings(destinationDir string, bearDB *db.DB) error {
	attachments, err := bearDB.AllAttachments()
	if err != nil {
		return errors.WithStack(err)
	}

	exportable := make([][]string, len(attachments)+1)

	exportable[0] = []string{"Note ID", "Note Title", "Attachment Path"}

	for i, a := range attachments {
		filepath := buildAttachmentFilename(destinationDir, a)
		exportable[i+1] = []string{a.NoteSHA, a.NoteTitle, filepath}
	}

	mappingFile, err := os.Create(path.Join(destinationDir, "Attachments.csv"))
	if err != nil {
		return errors.WithStack(err)
	}

	err = csv.NewWriter(mappingFile).WriteAll(exportable)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
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

func buildAttachmentFilename(directory string, attachment *db.Attachment) string {
	var dir string

	switch strings.ToLower(path.Ext(attachment.Filename)) {
	case ".jpeg", ".jpg", ".png", ".gif", ".tiff", ".tif", ".heic", ".heif":
		dir = "Note Images"
	default:
		dir = "Note Files"
	}

	return path.Join("Local Files", dir, attachment.FolderUUID, attachment.Filename)
}
