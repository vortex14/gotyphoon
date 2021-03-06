package interfaces

type MetaDataInterface interface {
	IsRequired() bool
	SetRequired()

	GetName() string
	SetName(name string)

	GetLabel()string
	SetLabel(label string)

	GetDescription() string
	SetDescription(description string)

	GetPath()string
	SetPath(path string)
}
