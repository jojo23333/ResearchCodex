package textutil

import (
	"strings"
	"unicode"
)

// Slugify converts a title into a filesystem-friendly slug.
func Slugify(input string) string {
	var b strings.Builder
	prevUnderscore := false

	for _, r := range strings.ToLower(strings.TrimSpace(input)) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			prevUnderscore = false
			continue
		}
		if r == ' ' || r == '-' || r == '_' {
			if !prevUnderscore {
				b.WriteRune('_')
				prevUnderscore = true
			}
			continue
		}
		// Skip everything else
	}

	res := b.String()
	res = strings.Trim(res, "_")
	if res == "" {
		return "idea"
	}
	return res
}
