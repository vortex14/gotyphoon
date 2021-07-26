package interfaces

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	FETCHER 				=	"Fetcher"
	PROCESSOR 				= 	"Processor"
	TRANSPORTER 			= 	"ResultTransporter"
	SCHEDULER 				= 	"Scheduler"
	DONOR 					= 	"Donor"

	DEFERRED				= 	"Deferred"
	PRIORITY				=	"Priority"
	RETRIES					=	"Retries"
	EXCEPTIONS				=	"Exceptions"

	PROCESSOR2PRIORITY 		= 	"ProcessorPriority"
	PROCESSOR2EXCEPTIONS	=	"ProcessorExceptions"

	FETCHER2PRIORITY		=	"FetcherPriority"
	FETCHER2EXCEPTIONS		=	"FetcherExceptions"

	SCHEDULER2PRIORITY		=	"SchedulerPriority"
	SCHEDULER2EXCEPTIONS	=	"SchedulerExceptions"

	TRANSPORTER2PRIORITY	=	"ResultTransporterPriority"

)


func (c *ConfigProject) SetConfigName(name string) {
	c.configFile = name
}

func (c *ConfigProject) GetConfigName() string {
	return c.configFile
}

func (c *ConfigProject) SetConfigPath(path string) {
	c.configPath = path
}

func (c *ConfigProject) GetConfigPath() string {
	return c.configPath
}


func (c *ConfigProject) GetComponentPort(name string) int {
	var port int
	switch name {
	case "donor":
		component := c.TyComponents.Donor
		port = component.Port
	case "fetcher":
		component := c.TyComponents.Fetcher
		port = component.Port
	case "processor":
		component := c.TyComponents.Processor
		port = component.Port
	case "result_transporter":
		component := c.TyComponents.ResultTransporter
		port = component.Port
	case "scheduler":
		component := c.TyComponents.Scheduler
		port = component.Port
	}

	return port
}


func (c *ConfigProject) GetConcurrent(component string, name string) int {
	return reflect.ValueOf(c.TyComponents).
		FieldByName(component).
		FieldByName("Queues").
		FieldByName(name).
		Interface().(Queue).Concurrent
}

func (c *ConfigProject) GetTopic(component string, name string, postFix string) string {

	settings := reflect.ValueOf(c.TyComponents).
		FieldByName(component).
		FieldByName("Queues").
		FieldByName(name).Interface().(Queue)

	projectTopicName := fmt.Sprintf("%s_%s", c.ProjectName, settings.Topic)

	if name == PRIORITY ||
		name == PROCESSOR2PRIORITY ||
		name == FETCHER2PRIORITY ||
		name == SCHEDULER2PRIORITY ||
		name == TRANSPORTER2PRIORITY {
		projectTopicName += postFix
	}

	if c.Debug {
		projectTopicName += "_debug"
	}

	if !settings.Share {
		projectTopicName = fmt.Sprintf("%s_%s", strings.ToLower(component), projectTopicName)
	}

	return projectTopicName
}



