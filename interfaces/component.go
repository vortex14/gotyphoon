package interfaces

type Component interface {
	goPromise
	GetName() string
	Stop(project Project)
	Start(project Project)
	Close(project Project)
	CheckDirectory(required []string, pathComponent string) bool
}
