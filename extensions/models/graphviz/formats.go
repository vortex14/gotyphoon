package graphviz

import "fmt"

func FormatBottomSpace(label string) string {
	return fmt.Sprintf("\n\n\n %s", label)
}

func FormatSpace(label string) string {
	return fmt.Sprintf("\n  %s  &nbsp;\n\n", label)
}