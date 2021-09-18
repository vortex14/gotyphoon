package interfaces

//import "github.com/goccy/go-graphviz/cgraph"
//
//type GraphOptions struct {
//	Name            string
//	IsCluster       bool
//	FontColor       string
//	Label           string
//	BackgroundColor string
//	PrefixNodeName  string
//	Style           string
//	Layout          string
//}
//
//type EdgeInterface interface {
//	Init() EdgeInterface
//	SetGraph(graph *cgraph.Graph) EdgeInterface
//
//	SetLabel(label string)
//	SetStyle(style string)
//	SetColor(color string)
//
//	SetArrowSize(size float64)
//	SetArrowHead(head string)
//	SetArrowTail(tail string)
//
//	SetHeadLabel(head string)
//	SetTailLabel(tail string)
//}
//
//type EdgeOptions struct {
//	Name string
//
//	NodeA string // from
//	NodeB string // to
//
//	Label string
//	LabelH string // label head
//	LabelT string // label tail
//
//	Color string
//	Style string
//
//	ArrowH string // default normal
//	ArrowT string
//	ArrowS float64 // arrow size 1 as default.
//}
//
//type NodeOptions struct {
//	Name string
//	Label string
//	Shape string
//	Style string
//	BackgroundColor string
//
//	*EdgeOptions
//}
//
//type NodeInterface interface {
//	Init() NodeInterface
//	SetParent(parentGraph *cgraph.Graph) NodeInterface
//	SetLabel(label string) NodeInterface
//	Get() *cgraph.Node
//	SetStyle(style string)
//	SetColor(color string)
//}
//
//type GraphInterface interface {
//	Init() GraphInterface
//	SetLayout(layout string)
//
//	Render(format string) []byte
//
//	AddNode(options *NodeOptions) GraphInterface
//	AddSubGraph(options *GraphOptions) GraphInterface
//	PostInit() GraphInterface
//
//	GetNodes() map[string]NodeInterface
//	SetNodes(nodes map[string]NodeInterface)
//
//	SetEdges(edges map[string]EdgeInterface)
//	GetEdges() map[string]EdgeInterface
//	UpdateEdge(options *EdgeOptions) GraphInterface
//	BuildEdges(nodesA []string, nodesB []string) GraphInterface
//}