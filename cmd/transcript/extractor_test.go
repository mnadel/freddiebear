package transcript

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMarkdownOneSection(t *testing.T) {
	markdown := `# Title
# Test Note
## First Section
#not/test
First sec text

## Second Section
#test
Second sec text
Even more
More second sec text
### Subsection
Sub text

## Third Section
#no/test
Third sec text
`

	extracted := `Second sec text
Even more
More second sec text
### Subsection
Sub text
`
	ex := NewTagExtractor([]byte(markdown), "#test")
	n := ex.ExtractTaggedNotes()

	assert.Equal(t, extracted, string(n))
}

func TestParseMarkdownTwoSections(t *testing.T) {
	markdown := `# Title
# Test Note
## First Section
#not/test
First sec text

## Second Section
#test
Second sec text
Even more
More second sec text
### Subsection
Sub text

## Third Section
#no/test
Third sec text

## Fourth Section
#test
More stuff to extract!
`

	extracted := `Second sec text
Even more
More second sec text
### Subsection
Sub text
More stuff to extract
!
`
	ex := NewTagExtractor([]byte(markdown), "#test")
	n := ex.ExtractTaggedNotes()

	assert.Equal(t, extracted, string(n))
}
