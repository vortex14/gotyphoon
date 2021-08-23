package services

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/vortex14/gotyphoon/integrations/mongo"
	"github.com/vortex14/gotyphoon/integrations/nsq"
	"github.com/vortex14/gotyphoon/integrations/redis"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/utils"
)

type Collections struct {
	Mongo map[string] *mongo.Service
	Redis map[string] *redis.Service
	Nsq *nsq.Service
}

type Services struct {
	Collections *Collections
	Project     interfaces.Project
	Config      *interfaces.ConfigProject
	Options 	interfaces.TyphoonIntegrationsOptions
}

func (s *Services) RunNSQ()  {
	s.initCollection()
	s.Collections.Nsq.RunNSQ()
}

func (s *Services) StopNSQ()  {
	println("Stop NSQ ...")
	s.Collections.Nsq.StopNSQ()
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
			Nsq: &nsq.Service{
				Project: s.Project,
				Options: s.Options.NSQ,
			},
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
