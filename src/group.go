package typhoon

type GroupLabel struct {
	Kind string
	Version string
}

type Group struct {
	Name string
	Description string
	Typhoon GroupLabel
	Servers []*Server
	Clouds []*Cloud
}

