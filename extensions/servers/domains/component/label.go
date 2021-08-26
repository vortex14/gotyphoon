package component

import (
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/interfaces/server"
)

func getDescription(componentName string) string {

	var description string

	switch componentName {
	case interfaces.SCHEDULER:
		description = SchedulerDescriptionComponent
	case interfaces.FETCHER:
		description = FetcherDescriptionComponent
	case interfaces.PROCESSOR:
		description = ProcessorDescriptionComponent
	case interfaces.TRANSPORTER:
		description = TransporterDescriptionComponent
	case interfaces.DONOR:
		description = DonorDescriptionComponent
	}

	return description
}

func ConstructorLabelOptions(componentName string) *server.BaseServerLabel {
	return &server.BaseServerLabel{
		Name: componentName,
		Description: getDescription(componentName),
	}
}
