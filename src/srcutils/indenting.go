package srcutils

import (
	"strings"
)

// @TODO: this isn't very pretty ...
func TrimConsistentIndenting(lines []string) []string {
	prefix := ""
	nPrefixes := 0
	prefixes := []string{" ", "\t"}
	for _, p := range prefixes {
		allLinesHavePrefix := true
		for _, l := range lines {
			// disregard the indenting of blank lines
			if l == "" {
				continue
			}

			if !strings.HasPrefix(l, p) {
				allLinesHavePrefix = false
				break
			}
		}

		if allLinesHavePrefix {
			prefix = p

			min := 100 // some way too high default, we'll min(min, ..) against this
			for _, l := range lines {
				// disregard the indenting of blank lines
				if l == "" {
					continue
				}

				i := 0
				for ; strings.HasPrefix(l, strings.Repeat(prefix, i+1)); i++ {
					// -
				}

				if i < min {
					min = i
				}
			}

			nPrefixes = min

			break
		}
	}

	if nPrefixes > 0 {
		trimPrefix := strings.Repeat(prefix, nPrefixes)

		for k, l := range lines {
			lines[k] = strings.TrimPrefix(l, trimPrefix)
		}
	}

	return lines
}
