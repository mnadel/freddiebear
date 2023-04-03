package alfred

import (
	"fmt"
	"strings"

	"github.com/mnadel/freddiebear/db"
	"github.com/mnadel/freddiebear/ext"
)

type Source = *db.Result
type Target = *db.Result

func AlfredBacklinkXML(matches map[Target]Source) string {
	builder := strings.Builder{}

	builder.WriteString(`<?xml version="1.0" encoding="utf-8"?>`)
	builder.WriteString(`<items>`)

	if len(matches) == 0 {
		builder.WriteString(`<item valid="no"><title>No backlinks found</title></item>`)
	} else {
		for target, source := range matches {
			source.Title = fmt.Sprintf("%s → %s", source.Title, target.Title)

			builder.WriteString(`<item valid="yes">`)
			builder.WriteString(`<title>`)
			builder.WriteString(source.TitleCase())
			builder.WriteString(`</title>`)

			builder.WriteString(`<subtitle>`)
			builder.WriteString(strings.Join(source.UniqueTags(), ", "))
			builder.WriteString(`</subtitle>`)

			builder.WriteString(`<arg>`)
			builder.WriteString(source.ID)
			builder.WriteString(`</arg>`)
			builder.WriteString(`</item>`)
		}
	}

	builder.WriteString(`</items>`)

	return builder.String()
}

func AlfredOpenXML(results db.Results, optShowTags bool) string {
	builder := strings.Builder{}

	builder.WriteString(`<?xml version="1.0" encoding="utf-8"?>`)
	builder.WriteString(`<items>`)

	for _, item := range results {
		builder.WriteString(`<item valid="yes">`)
		builder.WriteString(`<title>`)
		builder.WriteString(item.TitleCase())
		builder.WriteString(`</title>`)

		if !optShowTags {
			builder.WriteString(`<subtitle>Open note</subtitle>`)
		} else {
			builder.WriteString(`<subtitle>`)
			builder.WriteString(strings.Join(item.UniqueTags(), ", "))
			builder.WriteString(`</subtitle>`)
		}

		builder.WriteString(`<arg>`)
		builder.WriteString(item.ID)
		builder.WriteString(`</arg>`)
		builder.WriteString(`</item>`)
	}

	builder.WriteString(`</items>`)

	return builder.String()
}

func AlfredCreateXML(searchTerm string) string {
	builder := strings.Builder{}
	result := db.Result{
		Title: searchTerm,
	}
	title := result.TitleCase()

	builder.WriteString(`<?xml version="1.0" encoding="utf-8"?>`)
	builder.WriteString(`<items>`)

	builder.WriteString(`<item valid="yes">`)
	builder.WriteString(`<subtitle>Create note</subtitle>`)
	builder.WriteString(`<title>`)
	builder.WriteString(title)
	builder.WriteString(`</title>`)
	builder.WriteString(`<arg>`)
	ext.WriteKeyValue(&builder, `create`, title)
	builder.WriteString(`</arg>`)
	builder.WriteString(`</item>`)

	builder.WriteString(`</items>`)

	return builder.String()
}
