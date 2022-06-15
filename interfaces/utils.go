package interfaces

type FileSystem interface {
	GetDataFromDirectory(path string) MapFileObjects
	IsExistDir(path string) bool
	GetPath() string
}

type Utils interface {
	GoRunTemplate(goTemplate *GoTemplate) bool
	ParseLog(object *FileObject) error
	GetGoTemplate(object *FileObject) (error, string)
}
