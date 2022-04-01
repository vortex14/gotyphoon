package interfaces

type Component interface {
	goPromise
	GetName() string
	Start(project Project)
	Close(project Project)
	InitConsumers(project Project)
	IsActive() bool
	//CheckDirectory(required []string, pathComponent string) bool
}
