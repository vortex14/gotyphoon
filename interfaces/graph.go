package interfaces

type GraphOptions struct {
	Name            string
	IsCluster       bool
	FontColor       string
	Label           string
	BackgroundColor string
	PrefixNodeName  string
	Style           string
}

type EdgeOptions struct {
	Name string

	NodeA string // from
	NodeB string // to

	Label string
	LabelH string // label head
	LabelT string // label tail

	Color string
	Style string

	ArrowH string // default normal
	ArrowT string
	ArrowS float64 // arrow size 1 as default.
}

type NodeOptions struct {
	Name string
	Shape string

	*EdgeOptions
}

type GraphInterface interface {
	AddNode(options *NodeOptions) GraphInterface
	AddSubGraph(options *GraphOptions) GraphInterface
	Render(format string) []byte
	UpdateEdge(options *EdgeOptions) GraphInterface
	Init() GraphInterface
}