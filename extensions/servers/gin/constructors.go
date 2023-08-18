package gin

import (
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/interfaces"
)

func ConstructorCreateBaseLocalhostServer(name, description string, port int) *TyphoonGinServer {

	loggerOpt, tracingOpt, _ := ConstructorLocalhostOptions(name)

	return &TyphoonGinServer{
		TyphoonServer: &forms.TyphoonServer{
			IsDebug: true,
			Port:    port,
			Level:   interfaces.DEBUG,
			MetaInfo: &label.MetaInfo{
				Name:        name,
				Description: description,
			},
			TracingOptions: tracingOpt,
			LoggerOptions:  loggerOpt,
		},
	}
}
