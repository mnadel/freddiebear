package util

import (
	"log"
	"strings"
	"unicode"

	"github.com/pkg/errors"
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
	if len(sentence) == 0 || unicode.IsUpper(rune(sentence[0])) {
		return sentence
	}

	builder := strings.Builder{}

	for i, ch := range sentence {
		if i == 0 || sentence[i-1] == ' ' {
			builder.WriteRune(unicode.ToUpper(ch))
		} else {
			builder.WriteRune(ch)
		}
	}

	return builder.String()
}

// ToSafeString returns an Alfred-safe string
func ToSafeString(s string) string {
	return strings.ReplaceAll(s, "&", "&amp;")
}

// UniqueSet takes a string array and removes its duplicates
func UniqueSet(strs []string) []string {
	uniq := make(map[string]bool)
	uniques := make([]string, 0)

	for _, s := range strs {
		if !uniq[s] {
			uniq[s] = true
			uniques = append(uniques, s)
		}
	}

	return uniques
}

// MustString takes a (string, error) and either aborts the program
// if error isn't nil, or returns the string
func MustString(s string, e error) string {
	if e != nil {
		log.Fatalln(errors.WithStack(e))
	}

	return s
}
