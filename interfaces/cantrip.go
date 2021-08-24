package interfaces


type CantripInterface interface {
	Conjure(project Project) error
	GetName() string
}