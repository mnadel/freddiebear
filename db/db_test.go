package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniqueTags(t *testing.T) {
	r := Result{Tags: "fred,fred/bear,readings,work,work/coffee,work/coffee/africa"}
	unique := r.UniqueTags()

	assert.Equal(t, 3, len(unique), unique)

	assert.Contains(t, unique, "readings")
	assert.Contains(t, unique, "fred/bear")
	assert.Contains(t, unique, "work/coffee/africa")
}

func TestUniqueTagsSimple(t *testing.T) {
	r := Result{Tags: "fred,fred/bear"}
	unique := r.UniqueTags()

	assert.Equal(t, 1, len(unique), unique)

	assert.Contains(t, unique, "fred/bear")
}
