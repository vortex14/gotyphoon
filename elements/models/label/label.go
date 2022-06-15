package label

import "strings"

type MetaInfo struct {
	Path string
	Name string
	Label string
	Required bool
	Description string
}

func (l *MetaInfo) GetName() string {
	return l.Name
}

func (l *MetaInfo) SetLabel(label string)  {
	l.Label = label
}

func (l *MetaInfo) GetLabel() string {
	return l.Label
}

func (l *MetaInfo) GetDescription() string {
	return l.Description
}

func (l *MetaInfo) IsRequired() bool {
	return l.Required
}

func (l *MetaInfo) SetName(name string)  {
	l.Name = name
}

func (l *MetaInfo) SetDescription(description string)  {
	l.Description = description
}

func (l *MetaInfo) GetPath() string {
	return strings.ToLower(l.Path)
}

func (l *MetaInfo) SetPath(path string) {
	l.Path = path
}

func (l *MetaInfo) SetRequired()  {
	l.Required = true
}
