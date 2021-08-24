package interfaces



type DockerInterface interface {
	BuildImage()
	ListContainers()
	ProjectBuild()
	RunComponent(component string) error
}
