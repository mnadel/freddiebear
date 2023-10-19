package transcript

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mnadel/freddiebear/db"
)

func TestExtractTaggedNote(t *testing.T) {
	aNote := `
# Title

## Heading1
#tag
text1
text2

## Heading2
#tag/different
text3
text4

## Heading3
#tag
text5
text6
	`
	parts := extractTaggedNote(&db.Record{
		Text: aNote,
	}, "tag")
	
	assert.Equal(t, 2, len(parts))
	assert.Equal(t, "text1\ntext2", strings.TrimSpace(parts[0]))
	assert.Equal(t, "text5\ntext6", strings.TrimSpace(parts[1]))
}

