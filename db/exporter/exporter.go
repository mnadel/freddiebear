package exporter

import (
	"crypto/md5"
	"fmt"
	"io/fs"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
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
	files, err := ListFiles(directory)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	filenames := make(map[SHA]Filename)
	re := regexp.MustCompile(FilenameRegex)

	for _, file := range files {
		parts := re.FindStringSubmatch(file)
		if len(parts) == 2 {
			filenames[SHA(parts[1])] = Filename(file)
		}
	}

	return &Exporter{filenames, directory}, nil
}

// Archive will move archived notes to trashDirectory
func (e *Exporter) Archive(records []*db.Record, trashDirectory string) error {
	// create lookup table for current records
	currSHAs := make(map[SHA]bool)
	for _, rec := range records {
		currSHAs[SHA(rec.SHA)] = true
	}

	// iterate over the list of exported notes
	for sha, file := range e.mapping {
		// and if it's not in the list of currents
		if ok := currSHAs[sha]; !ok {
			// then move to the trash directory
			log.Println("archiving", string(file))
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

func ListFiles(directory string) ([]string, error) {
	var files []string

	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if ignore, err := ignorePath(directory, path, d.IsDir()); err != nil {
			return err
		} else if !ignore {
			files = append(files, d.Name())
		}

		return nil
	})

	return files, err
}

func ignorePath(cwd, path string, isDir bool) (bool, error) {
	if isDir {
		return true, nil
	}

	cwdDir, err := filepath.Abs(cwd)
	if err != nil {
		return false, err
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return false, err
	}

	remaining := strings.TrimPrefix(absPath, cwdDir)
	if len(remaining) > 0 {
		remaining = remaining[1:]
	}
	return strings.Contains(remaining, PathSep), nil
}
