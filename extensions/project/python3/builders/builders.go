package builders

type Builders interface {
	BuilderV11()
}

type BuildOptions struct {
	Component string
	Type string
	Meta interface{}

}


type ProjectBuilder interface {
	Build()
}
