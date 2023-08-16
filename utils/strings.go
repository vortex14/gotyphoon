package utils

import (
	"strings"
	"unicode"
)

func IsStrContain(source string, args ...string) bool {
	status := false

	for _, term := range args {
		if strings.Contains(source, term) {
			status = true
			break
		}
	}

	return status
}

func ClearStrTabAndN(source string) string {
	return strings.ReplaceAll(strings.ReplaceAll(source, "\n", ""), "\t", "")
}

func IsFirstUpLetter(str string) bool {
	status := false
	for _, v := range []rune(str) {
		if unicode.IsUpper(v) {
			status = true
			break
		}
	}
	return status
}
