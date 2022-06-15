package main

import (
	"github.com/vortex14/gotyphoon/elements/models/label"
	. "github.com/vortex14/gotyphoon/extensions/agents/events-metric-v1"
)

func main() {
	agent := AgentMetric{
		MetaInfo:          &label.MetaInfo{
			Name: "AgentMetric",
			Description: "AgentMetric Server",
		},
		ServerDescription: "Agent metric server",
		ServerBasePath:    "/",
		ServerPort:        14000,
		ServerName:        "Agent metric",
		Databases:         []string{"events"},
		MongoHost:         "localhost",
		MongoPort:         27017,
		OutDb:             "events",
		OutCollection:     "tracks",
	}
	agent.Run()
}