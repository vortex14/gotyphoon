package nsq

import (
	"encoding/json"
	"os"
	"reflect"

	"github.com/fatih/color"
	"github.com/fatih/structs"
	"github.com/vortex14/gotyphoon/interfaces"
)



func (s *Service) initConfig()  {
	if s.Config == nil {
		//color.Red("%+v", s.Options)
		s.priorityConsumers = make(map[string]map[int]*Consumer)
		s.priorityProducers = make(map[string]map[int]*Producer)
		s.Config = s.Project.LoadConfig()
	}
}

func (s *Service) InitQueue(settings *interfaces.Queue)  {
	if settings.Readable && s.Options.EnabledConsumer {
		s.InitConsumer(settings)
	}

	if settings.Writable && s.Options.EnabledProducer {
		s.InitProducer(settings)
	}
}


func (s *Service) getSettings() <-chan *interfaces.Queue {
	ch := make(chan*interfaces.Queue)

	go func(ch chan *interfaces.Queue) {
		for _, component := range s.Project.GetSelectedComponent() {
			queueSettings := reflect.ValueOf(s.Config.TyComponents).
				FieldByName(component).
				FieldByName("Queues").Interface()
			queueSettingsMap := structs.Map(queueSettings)
			for groupName, settingsMap := range queueSettingsMap {
				sourceData, _ := json.Marshal(settingsMap)
				var qSettings interfaces.Queue
				err := json.Unmarshal(sourceData, &qSettings)
				if err != nil {
					color.Red("%s", err.Error())
					os.Exit(1)
				}
				qSettings.SetGroupName(groupName)
				qSettings.SetComponentName(component)
				ch <- &qSettings
			}
		}


		defer close(ch)
	}(ch)

	return ch
}



