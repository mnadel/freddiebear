package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniqueTags(t *testing.T) {
	orig := []string{"fred", "fred/bear", "readings", "work", "work/coffee", "work/coffee/africa"}
	unique := RemoveIntermediatePrefixes(orig, "/")

	assert.Equal(t, 3, len(unique), unique)

	assert.Contains(t, unique, "readings")
	assert.Contains(t, unique, "fred/bear")
	assert.Contains(t, unique, "work/coffee/africa")

	assert.Equal(t, 6, len(orig))
}
