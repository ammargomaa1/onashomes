package utils

import (
	"regexp"
	"strings"
)

var (
	slugNonAlphaNum = regexp.MustCompile(`[^a-z0-9\s-]`)
	slugWhitespace  = regexp.MustCompile(`[\s]+`)
	slugMultiDash   = regexp.MustCompile(`-+`)
)

// Slugify converts a string to a URL-friendly slug
func Slugify(input string) string {
	s := strings.ToLower(strings.TrimSpace(input))
	s = slugNonAlphaNum.ReplaceAllString(s, "")
	s = slugWhitespace.ReplaceAllString(s, "-")
	s = slugMultiDash.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}
