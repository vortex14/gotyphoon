package label

import "strings"

type MetaInfo struct {
	Path        string
	Name        string
	Label       string
	Required    bool
	Version     string
	Description string
	Summary     string
	Tags        []string
}

func (l *MetaInfo) GetName() string {
	return l.Name
}

func (l *MetaInfo) SetLabel(label string) {
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

func (l *MetaInfo) SetName(name string) {
	l.Name = name
}

func (l *MetaInfo) SetDescription(description string) {
	l.Description = description
}

func (l *MetaInfo) GetPath() string {
	return strings.ToLower(l.Path)
}

func (l *MetaInfo) SetPath(path string) {
	l.Path = path
}

func (l *MetaInfo) SetRequired() {
	l.Required = true
}

func (l *MetaInfo) GetTags() []string {
	return l.Tags
}

func (l *MetaInfo) SetTag(tag string) {
	l.Tags = append(l.Tags, tag)
}

func (l *MetaInfo) GetSummary() string {
	return l.Summary
}
