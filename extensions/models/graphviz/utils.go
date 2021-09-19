package graphviz
//
// /* ignore for building amd64-linux
import (
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

func GetStyle(style string) cgraph.GraphStyle {
	return cgraph.GraphStyle(style)
}

func GetExportFormat(format string) graphviz.Format {
	return graphviz.Format(format)
}
// */