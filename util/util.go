package util

import (
	"strings"
)

// RemoveIntermediatePrefixes removes the set of intermediate prefixes. If [a a/b c] is passed in,
// then [a/b c] is returned; `a` was removed because it's an intermediate prefix of `a/b` (given a separator of /).
func RemoveIntermediatePrefixes(strs []string, sep string) []string {
	// create a copy that we can mutate
	mut := strs

	prefix := strings.Builder{}

	// if any tag is a prefix of another, clear it out
	for i, s := range mut {
		if s == "" {
			continue
		}

		prefix.WriteString(s)
		prefix.WriteString(sep)
		pfx := prefix.String()

		for j, t := range mut {
			if i != j && mut[i] != "" && strings.HasPrefix(t, pfx) {
				mut[i] = ""
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
