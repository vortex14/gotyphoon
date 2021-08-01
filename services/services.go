package services

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/fatih/structs"
	"github.com/vortex14/gotyphoon/integrations/mongo"
	"github.com/vortex14/gotyphoon/integrations/nsq"
	"github.com/vortex14/gotyphoon/integrations/redis"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/utils"
	"os"
	"reflect"
	"strconv"
)

type Collections struct {
	Mongo map[string] *mongo.Service
	Redis map[string] *redis.Service
	Nsq *nsq.Service
}

type Services struct {
	Project     interfaces.Project
	Collections *Collections
	Config      *interfaces.ConfigProject
}

func (s *Services) getSettings() <-chan *interfaces.Queue {
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
				//color.Yellow("%+v", qSettings)
				ch <- &qSettings
			}
		}


		defer close(ch)
	}(ch)

	return ch
}

func (s *Services) RunNSQ()  {
	fmt.Println("Running connections to NSQ ...")
	nsqService := s.Collections.Nsq

	for qSettings := range s.getSettings() {
		group := qSettings.GetGroupName()
		if group == interfaces.PRIORITY ||
			group == interfaces.PROCESSOR2PRIORITY ||
			group == interfaces.TRANSPORTER2PRIORITY {
			color.Yellow("Init Priority queues ...")
			for _, i := range []int{1,2,3} {
				prioritySetting := &interfaces.Queue{
					Concurrent: qSettings.Concurrent,
					MsgTimeout: qSettings.MsgTimeout,
					Channel:    qSettings.Channel,
					Topic:      s.Config.GetTopic(qSettings.GetComponentName(), group, strconv.Itoa(i)),
					Share:      qSettings.Share,
					Writable:   qSettings.Writable,
					Readable:   qSettings.Readable,
				}
				prioritySetting.SetPriority(i)
				prioritySetting.SetGroupName(group)
				nsqService.InitQueue(prioritySetting)
			}

			continue
		}
		qSettings.Topic = s.Config.GetTopic(qSettings.GetComponentName(), group, "")
		nsqService.InitQueue(qSettings)
	}
}

func (s *Services) StopNSQ()  {
	s.Collections.Nsq.StopProducers()
	s.Collections.Nsq.StopConsumers()
}

func (s *Services) RunTestServices() {
	s.initServices()
	var tableData [][]string
	stackType := s.GetStackType()
	header := []string{"Name", "Group", "Type", "Status"}
	for group, service := range s.Collections.Mongo {
		status := strconv.FormatBool(service.Ping())
		tableData = append(tableData, []string{"Mongo",group, stackType ,status})
	}

	for group, service := range s.Collections.Redis {
		status := strconv.FormatBool(service.Ping())
		tableData = append(tableData, []string{"Redis",group, stackType ,status})
	}

	nsqStatus := strconv.FormatBool(s.Collections.Nsq.Ping())
	tableData = append(tableData, []string{"Nsq", "main", stackType, nsqStatus})
	u := utils.Utils{}
	u.RenderTableOutput(header, tableData)
}

func (s *Services) GetStackType() string {
	var stackType string
	if s.Config.Debug {
		stackType = "Debug"
	} else {
		stackType = "Production"
	}
	return stackType
}

func (s *Services) GetMongoStack() []interfaces.ServiceMongo {
	mongoStack := reflect.
		ValueOf(s.Config.Services.Mongo).
		FieldByName(s.GetStackType()).
		Interface().([]interfaces.ServiceMongo)
	return mongoStack
}

func (s *Services) GetRedisStack() []interfaces.ServiceRedis {
	redisStack := reflect.
		ValueOf(s.Config.Services.Redis).
		FieldByName(s.GetStackType()).
		Interface().([]interfaces.ServiceRedis)
	return redisStack
}

func (s *Services) initStack() {
	stackType := s.GetStackType()
	fmt.Println("Service stack Type: "+ stackType)
	mongoStack := s.GetMongoStack()
	redisStack := s.GetRedisStack()

	for _, service := range mongoStack {
		s.initMongoService(service)
	}

	for _, service := range redisStack {
		s.initRedisService(service)
	}
}

func (s *Services) initCollection()  {
	if s.Collections == nil {
		s.Config = s.Project.LoadConfig()
		s.Collections = &Collections{
			Mongo: map[string]*mongo.Service{},
			Redis: map[string]*redis.Service{},
			Nsq: &nsq.Service{Project: s.Project},
		}
	}
}

func (s *Services) initServices()  {
	if s.Collections == nil {
		s.initCollection()
		s.initStack()
	}
}

func (s *Services) initMongoService(service interfaces.ServiceMongo)  {
	mongoService := mongo.Service{Project: s.Project, Settings: service}
	s.Collections.Mongo[service.Name] = &mongoService
	mongoService.Init()
}

func (s *Services) initRedisService(service interfaces.ServiceRedis)  {
	redisService := redis.Service{Project: s.Project, Config: &service}
	s.Collections.Redis[service.Name] = &redisService
	redisService.Init()

}

func (s *Services) LoadMongoServices()  {
	s.initCollection()
	mongoStack := s.GetMongoStack()
	for _, service := range mongoStack {
		s.initMongoService(service)
	}
}

func (s *Services) LoadRedisServices()  {
	s.initCollection()
	redisStack := s.GetRedisStack()

	for _, service := range redisStack {
		s.initRedisService(service)
	}
}

func (s *Services) LoadProjectServices()  {
	fmt.Println("Load project services ...")
	s.initServices()
	//s.LoadMongoServices()
	//s.LoadRedisServices()
}
