package util

import (
	"strings"
	"unicode"
)

// RemoveIntermediatePrefixes removes the set of intermediate prefixes. If [a a/b c] is passed in,
// then [a/b c] is returned; `a` was removed because it's an intermediate prefix of `a/b` (given a separator of /).
func RemoveIntermediatePrefixes(strs []string, sep string) []string {
	// create a copy that we can mutate
	mut := strs

	prefix := strings.Builder{}

	for i, s := range mut {
		// if it's already been removed, then skip
		if s == "" {
			continue
		}

		prefix.WriteString(s)
		prefix.WriteString(sep)
		pfx := prefix.String()

		// remove if any other string is prefixed with this string
		for j, t := range mut {
			if i != j && strings.HasPrefix(t, pfx) {
				mut[i] = ""
				break
			}
		}

		prefix.Reset()
	}

	// collect non-empty entries
	collapsed := make([]string, 0)

	for _, s := range mut {
		if s != "" {
			collapsed = append(collapsed, s)
		}
	}

	return collapsed
}

// ToTitleCase returns the string with first char of each word uppercased.
// Unless the first letter already is, then it just returns the string.
func ToTitleCase(sentence string) string {
	if unicode.IsUpper(rune(sentence[0])) {
		return sentence
	}

	builder := strings.Builder{}

	for _, word := range strings.Split(sentence, " ") {
		if builder.Len() > 0 {
			builder.WriteString(" ")
		}
		builder.WriteRune(unicode.ToUpper(rune(word[0])))
		builder.WriteString(word[1:])
	}

	return builder.String()
}
