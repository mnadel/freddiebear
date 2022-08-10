package util

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveIntermediatePrefixes(t *testing.T) {
	orig := []string{"fred", "fred/bear", "readings", "work", "work/coffee", "work/coffee/africa"}
	unique := RemoveIntermediatePrefixes(orig, "/")

	assert.Equal(t, 3, len(unique), unique)

	assert.Contains(t, unique, "readings")
	assert.Contains(t, unique, "fred/bear")
	assert.Contains(t, unique, "work/coffee/africa")

	assert.Equal(t, 6, len(orig))
}

func TestToTitleCase(t *testing.T) {
	// single word
	assert.Equal(t, "Bobby", ToTitleCase("bobby"))

	// two words
	assert.Equal(t, "Bobby Tables", ToTitleCase("bobby tables"))

	// bunch of words
	assert.Equal(t, "My Name Is Bobby Tables", ToTitleCase("my name is bobby tables"))

	// word that also has caps
	assert.Equal(t, "QrstuVwX", ToTitleCase("qrstuVwX"))

	// preserve when first letter is already uppercase
	assert.Equal(t, "Bobby tables", ToTitleCase("Bobby tables"))
}

func TestToSafeString(t *testing.T) {
	// ampersand isn't safe
	assert.Equal(t, "a &amp; b", ToSafeString("a & b"))

	// frontslash is safe
	assert.Equal(t, "a / b", ToSafeString("a / b"))
}

func TestUniqueSet(t *testing.T) {
	test1 := []string{"a", "a", "b", "a"}
	set1 := UniqueSet(test1)

	assert.Equal(t, 2, len(set1), set1)
	assert.Contains(t, set1, "a")
	assert.Contains(t, set1, "b")
}

func BenchmarkRemoveIntermediatePrefixes(t *testing.B) {
	tests := [][]string{
		{"fred", "fred/bear", "readings", "work", "work/coffee", "work/coffee/africa"},
		{"a", "b", "c", "d"},
		{"a", "a/b", "a/b/c", "a/b/c/d"},
		{"a/b/c/d", "a/b/c", "a/b", "a"},
	}

	for i := 0; i < t.N; i++ {
		t := rand.Intn(len(tests))
		RemoveIntermediatePrefixes(tests[t], "/")
	}
}

func BenchmarkTitleCase(t *testing.B) {
	for i := 0; i < t.N; i++ {
		ToTitleCase("bobby tables")
	}
}
