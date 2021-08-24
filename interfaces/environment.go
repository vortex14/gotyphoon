package interfaces

import "github.com/vortex14/gotyphoon/environment"

type Environment interface {
	Load()
	Set()
	Get()
	GetSettings() (error, *environment.Settings)
}
