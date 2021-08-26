package server

type Action struct {
	Name string
	Methods []string
	Description  string
	Controller   Controller
	PyController Controller
}

type ActionInterface interface {
	AddMethod(name string) error
	MetaDataInterface
}
