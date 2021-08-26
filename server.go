package typhoon

import "github.com/vortex14/gotyphoon/services"

type ServerLabel struct {
	Kind string
	Version string
}

type Server struct {
	Name        string
	Description string
	Clusters    []*Cluster
	Typhoon     ServerLabel
	Services    *services.Services
}
