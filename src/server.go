package typhoon

type ServerLabel struct {
	Kind string
	Version string
}

type Server struct {
	Name string
	Description string
	Clusters []*Cluster
	Typhoon ServerLabel
}
