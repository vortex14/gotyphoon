package utils

import "strings"

func IsStrContain(source string, args ...string) bool {
	status := false

	for _, term := range args {
		if strings.Contains(source, term) { status = true; break }
	}

	return status
}
