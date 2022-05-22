package exporter

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/mnadel/freddiebear/db"
	"github.com/pkg/errors"
)

const (
	FilenameTemplate = "%s (%s).md"
	FilenameRegex    = `.*\s\((\w+)\)\.md$`
	PathSep          = string(os.PathSeparator)
)

type SHA string
type Filename string

type Exporter struct {
	mapping   map[SHA]Filename
	directory string
}

func NewExporter(directory string) (*Exporter, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	filenames := make(map[SHA]Filename)
	re := regexp.MustCompile(FilenameRegex)

	for _, file := range files {
		if !file.IsDir() {
			parts := re.FindStringSubmatch(file.Name())
			if len(parts) == 2 {
				filenames[SHA(parts[1])] = Filename(file.Name())
			}
		}
	}

	return &Exporter{filenames, directory}, nil
}

// Archive will move archived notes to trashDirectory
func (e *Exporter) Archive(records []*db.Record, trashDirectory string) error {
	currSHAs := make(map[SHA]bool)

	for _, rec := range records {
		currSHAs[SHA(rec.SHA)] = true
	}

	for sha, file := range e.mapping {
		if ok := currSHAs[sha]; !ok {
			newName := path.Join(trashDirectory, string(file))
			if err := os.Rename(string(file), newName); err != nil {
				return errors.WithStack(err)
			}
		}
	}

	return nil
}

// Returns true if the SHA and its new data differs from the previously-exported contents
func (e *Exporter) IsChanged(record *db.Record) (bool, error) {
	filename, ok := e.mapping[SHA(record.SHA)]
	if !ok {
		return true, nil
	}

	oldData, err := os.ReadFile(string(filename))
	if err != nil {
		return false, errors.WithStack(err)
	}

	oldSum := md5.Sum(oldData)
	newSum := md5.Sum([]byte(record.Text))

	return oldSum != newSum, nil
}

// Returns true if the SHA has been renamed, and if so, what the previous name was
func (e *Exporter) IsRenamed(record *db.Record) (bool, Filename) {
	f, ok := e.mapping[SHA(record.SHA)]
	if !ok {
		return false, ""
	}

	return string(f) != BuildFilename(record), f
}

func BuildFilename(record *db.Record) string {
	safeTitle := strings.ReplaceAll(record.Title, PathSep, url.QueryEscape(PathSep))

	return fmt.Sprintf(FilenameTemplate, safeTitle, record.SHA)
}
