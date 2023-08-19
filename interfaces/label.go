package interfaces

type MetaDataInterface interface {
	IsRequired() bool
	SetRequired()

	GetName() string
	SetName(name string)

	GetLabel() string
	SetLabel(label string)

	GetDescription() string
	GetSummary() string
	SetDescription(description string)

	SetTag(tag string)
	GetTags() []string

	GetPath() string
	SetPath(path string)
}
