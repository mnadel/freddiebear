package transcript

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

var newline = []byte("\n")

type TagExtractor struct {
	foundTag bool
	buf      bytes.Buffer
	tag      []byte
	source   []byte
}

func NewTagExtractor(source []byte, tag string) *TagExtractor {
	return &TagExtractor{
		foundTag: false,
		buf:      bytes.Buffer{},
		tag:      []byte(tag),
		source:   source,
	}
}

func (te *TagExtractor) ExtractTaggedNotes() []byte {
	doc := goldmark.New().Parser().Parse(text.NewReader(te.source))

	if optAst {
		ast.DumpHelper(doc, te.source, 0, nil, func(_ int) {})
		return nil
	}

	ast.Walk(doc, te.visit)

	return te.buf.Bytes()
}

func (te *TagExtractor) visit(node ast.Node, entering bool) (ast.WalkStatus, error) {
	headingLevel := 0

	if entering {
		debug("entering %s, text=%s", node.Kind(), node.Text(te.source))
	} else {
		debug("visiting %s", node.Kind())
	}

	if header, isHeader := node.(*ast.Heading); isHeader {
		headingLevel = header.Level
		if te.foundTag && headingLevel < 3 {
			te.foundTag = false
			debug("resetting found tag")
		}
	} else if para, isPara := node.(*ast.Paragraph); isPara {
		if bytes.Contains(para.Text(te.source), te.tag) {
			te.foundTag = true
			debug("found tag")
		}
	} else {
		return ast.WalkContinue, nil
	}

	if !entering && te.foundTag {
		te.buf.WriteString(te.nodeText(node, headingLevel))
	}

	return ast.WalkSkipChildren, nil
}

func (te *TagExtractor) nodeText(node ast.Node, headingLevel int) string {
	s := strings.Builder{}
	n := node.FirstChild()

	for {
		if n == nil {
			break
		}

		text := n.Text(te.source)

		if headingLevel > 0 {
			s.WriteString(mdHeading(headingLevel))
			s.WriteString(" ")
		}

		if !bytes.Equal(text, te.tag) {
			s.Write(text)
			s.Write(newline)
		}

		n = n.NextSibling()
	}

	return s.String()
}

func mdHeading(level int) string {
	s := strings.Builder{}

	for i := 0; i < level; i++ {
		s.WriteString("#")
	}

	return s.String()
}

func debug(msg string, args ...interface{}) {
	if optDebug {
		fmt.Printf(msg+"\n", args...)
	}
}

/*

Document {
    Heading {
        RawText: "Test Note"
        HasBlankPreviousLines: true
        Level: 1
        Text: "Test Note"
    }
    Heading {
        RawText: "First Section"
        HasBlankPreviousLines: false
        Level: 2
        Text: "First Section"
    }
    Paragraph {
        RawText: "#nottest
First sec text"
        HasBlankPreviousLines: false
        Text(SoftLineBreak): "#nottest"
        Text: "First sec text"
    }
    Heading {
        RawText: "Second Section"
        HasBlankPreviousLines: true
        Level: 2
        Text: "Second Section"
    }
    Paragraph {
        RawText: "#test
Second sec text
Even more
More second sec text"
        HasBlankPreviousLines: false
        Text(SoftLineBreak): "#test"
        Text(SoftLineBreak): "Second sec text"
        Text(SoftLineBreak): "Even more"
        Text: "More second sec text"
    }
    Heading {
        RawText: "Subsection"
        HasBlankPreviousLines: false
        Level: 3
        Text: "Subsection"
    }
    Paragraph {
        RawText: "Sub text"
        HasBlankPreviousLines: false
        Text: "Sub text"
    }
    Heading {
        RawText: "Third Section"
        HasBlankPreviousLines: true
        Level: 2
        Text: "Third Section"
    }
    Paragraph {
        RawText: "#testno
Third sec text"
        HasBlankPreviousLines: false
        Text(SoftLineBreak): "#testno"
        Text: "Third sec text"
    }
}

Document {
    Heading {
        RawText: "Test Note"
        HasBlankPreviousLines: true
        Level: 1
        Text: "Test Note"
    }
    Heading {
        RawText: "First Section"
        HasBlankPreviousLines: false
        Level: 2
        Text: "First Section"
    }
    Paragraph {
        RawText: "First sec text"
        HasBlankPreviousLines: false
        Text: "First sec text"
    }
    Heading {
        RawText: "Second Section"
        HasBlankPreviousLines: true
        Level: 2
        Text: "Second Section"
    }
    Paragraph {
        RawText: "#test
Second sec text
More second sec text"
        HasBlankPreviousLines: false
        Text(SoftLineBreak): "#test"
        Text(SoftLineBreak): "Second sec text"
        Text: "More second sec text"
    }
    Heading {
        RawText: "Third Section"
        HasBlankPreviousLines: true
        Level: 2
        Text: "Third Section"
    }
    Paragraph {
        RawText: "Third sec text"
        HasBlankPreviousLines: false
        Text: "Third sec text"
    }
}

*/
