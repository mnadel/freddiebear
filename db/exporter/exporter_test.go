package exporter

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/mnadel/freddiebear/db"
	"github.com/stretchr/testify/assert"
)

func TestDetectRename(t *testing.T) {
	exp, err := NewExporter(".")
	assert.NoError(t, err)

	exp.mapping[SHA("abc123")] = Filename("previous_title (abc123).md")
	record := &db.Record{
		SHA:   "abc123",
		Title: "new_title",
		Text:  "",
	}

	changed, oldName := exp.IsRenamed(record)

	assert.True(t, changed)
	assert.Equal(t, "previous_title (abc123).md", string(oldName))
	assert.Equal(t, "new_title (abc123).md", BuildFilename(record))
}

func TestDetectNoRename(t *testing.T) {
	exp, err := NewExporter(".")
	assert.NoError(t, err)

	exp.mapping[SHA("abc123")] = Filename("previous_title (abc123).md")
	record := &db.Record{
		SHA:   "abc123",
		Title: "previous_title",
		Text:  "",
	}

	changed, oldName := exp.IsRenamed(record)

	assert.False(t, changed)
	assert.Equal(t, "previous_title (abc123).md", string(oldName))
	assert.Equal(t, "previous_title (abc123).md", BuildFilename(record))
}

func TestDetectChange(t *testing.T) {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "freddiebear")
	assert.NoError(t, err)

	tmpFile, err := ioutil.TempFile(tmpDir, "freddiebear-")
	assert.NoError(t, err)

	defer os.Remove(tmpFile.Name())

	exp, err := NewExporter(tmpDir)
	assert.NoError(t, err)

	exp.mapping[SHA("abc123")] = Filename(tmpFile.Name())
	assert.NoError(t, os.WriteFile(tmpFile.Name(), []byte("original content"), 0644))

	record := &db.Record{
		SHA:   "abc123",
		Title: "new_title",
		Text:  "new and improved updated content",
	}

	changed, err := exp.IsChanged(record)
	assert.NoError(t, err)
	assert.True(t, changed)
}

func TestDetectNoChange(t *testing.T) {
	tmpDir, err := os.MkdirTemp(os.TempDir(), "freddiebear")
	assert.NoError(t, err)

	tmpFile, err := ioutil.TempFile(tmpDir, "freddiebear-")
	assert.NoError(t, err)

	defer os.Remove(tmpFile.Name())

	exp, err := NewExporter(tmpDir)
	assert.NoError(t, err)

	exp.mapping[SHA("abc123")] = Filename(tmpFile.Name())
	assert.NoError(t, os.WriteFile(tmpFile.Name(), []byte("original content"), 0644))

	record := &db.Record{
		SHA:   "abc123",
		Title: "new_title",
		Text:  "original content",
	}

	changed, err := exp.IsChanged(record)
	assert.NoError(t, err)
	assert.False(t, changed)
}

func TestListFiles(t *testing.T) {
	dir, err := os.Getwd()
	assert.NoError(t, err)

	oldMethodFiles, err := ioutil.ReadDir(dir)
	assert.NoError(t, err)

	newMethodFiles, err := ListFiles(dir)
	assert.NoError(t, err)

	assert.Equal(t, len(oldMethodFiles), len(newMethodFiles))

	oldNames := make([]string, len(oldMethodFiles))
	for _, e := range oldMethodFiles {
		oldNames = append(oldNames, e.Name())
	}

	newNames := make([]string, len(newMethodFiles))
	for _, e := range newMethodFiles {
		newNames = append(newNames, e.Name())
	}

	assert.ElementsMatch(t, oldNames, newNames)
}
