package util

import (
	"strings"
)

// RemoveIntermediatePrefixes removes the set of intermediate prefixes. If [a a/b c] is passed in,
// then [a/b c] is returned; `a` was removed because it's an intermediate prefix of `a/b` (given a separator of /).
func RemoveIntermediatePrefixes(strs []string, sep string) []string {
	mut := strs

	// if any tag is a prefix of another, wipe it out
	prefix := strings.Builder{}

	for i, tag := range mut {
		for j, t := range mut {
			prefix.WriteString(tag)
			prefix.WriteString(sep)

			if i != j && mut[i] != "" && strings.HasPrefix(t, prefix.String()) {
				mut[i] = ""
			}

			prefix.Reset()
		}
	}

	// collect non-empty entries
	collapsed := make([]string, 0)

	for _, tag := range mut {
		if tag != "" {
			collapsed = append(collapsed, tag)
		}
	}

	return collapsed
}
