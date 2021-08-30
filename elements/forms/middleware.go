package forms

import "github.com/vortex14/gotyphoon/interfaces"

type Middleware struct {
	*interfaces.BaseLabel
}

func (m *Middleware) GetName() string {
	return m.Name
}

func (m *Middleware) GetDescription() string {
	return m.Description
}