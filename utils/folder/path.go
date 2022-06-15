package folder

import "strings"

// ValidRelPath check for path traversal and correct forward slashes
func ValidRelPath(p string) bool {
	if p == "" ||
		strings.Contains(p, `\`) ||
		strings.HasPrefix(p, "/") ||
		strings.Contains(p, "../") {
		return false
	}
	return true
}
