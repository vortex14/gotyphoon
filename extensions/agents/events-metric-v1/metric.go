package events_metric_v1

import (
	"context"
	"time"

	Gin "github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/awaitabler"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/elements/models/singleton"
	"github.com/vortex14/gotyphoon/extensions/servers/gin"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/controllers/graph"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/controllers/ping"
	GinMiddlewares "github.com/vortex14/gotyphoon/extensions/servers/gin/middlewares"
	"github.com/vortex14/gotyphoon/integrations/mongo"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/interfaces/server"
	"github.com/vortex14/gotyphoon/log"

	GinExtension "github.com/vortex14/gotyphoon/extensions/servers/gin"
	Mongo "go.mongodb.org/mongo-driver/mongo"
)

const (
	TRACKER = "tracker"
)


type Track struct {
	Event        string      `json:"event" binding:"required"`
	RegisteredAt time.Time   `json:"registered_at" binding:"required"`
	Payload      interface{} `json:"payload"`
}

type AgentMetric struct {
	singleton.Singleton
	awaitabler.Object
	*label.MetaInfo

	ServerDescription string
	ServerBasePath    string
	ServerPort        int
	ServerName        string
	Databases         []string
	MongoHost         string
	MongoPort         int

	OutDb         string
	OutCollection string
	outCollection *Mongo.Collection

	mongoService *mongo.Service
	server       *gin.TyphoonGinServer
	LOG          interfaces.LoggerInterface
}

func (m *AgentMetric) Run()  {
	m.Construct(func() {
		m.LOG = log.New(log.D{"agent": "metric", "created_at": time.Now().String()})
		m.LOG.Info("init !")

		m.mongoService = mongo.CreateMongoServiceWithoutAuth(m.Name, m.MongoHost, m.MongoPort, m.Databases)
		m.outCollection = m.mongoService.GetMongoCollection(m.OutDb, m.OutCollection)

		m.createServer()

		errS := m.server.Run()
		if errS != nil {
			return
		}

	})
}

func (m *AgentMetric) createServer()  {
	m.server = gin.ConstructorCreateBaseLocalhostServer(m.ServerName, m.ServerDescription, m.ServerPort)
	m.server.
		Init().
		AddResource(
			&forms.Resource{
				MetaInfo: &label.MetaInfo{
					Name:        m.ServerName,
					Path:        m.ServerBasePath,
					Description: m.ServerDescription,
				},
				Actions: map[string]interfaces.ActionInterface{
					ping.PATH:  ping.Controller,
					graph.PATH: graph.Controller,

					TRACKER: &GinExtension.Action{
						Action: &forms.Action{
							Middlewares: []interfaces.MiddlewareInterface{
								GinMiddlewares.ConstructorCorsMiddleware(server.GetAllAllowedCors()),
							},
							MetaInfo: &label.MetaInfo{
								Name:        TRACKER,
								Path:        TRACKER,
								Description: "listen new data track",
							},
							Methods: []string{interfaces.POST, interfaces.OPTIONS, interfaces.GET},
						},
						GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {

							var track Track
							err := ctx.ShouldBindJSON(&track)
							if err != nil {
								ctx.JSON(400,
									Gin.H{"status": false, "msg": "event, registered_at required json fields!"})
								return
							}

							filter :=  m.mongoService.GetFilterOptions("event", track.Event)
							update := m.mongoService.GetIncOptions(bson.D{{"count", 1}})

							now := time.Now()
							update = m.mongoService.GetIncUpdateOptions(
								bson.D{{"count", 1}},
								bson.D{
									{"registered_at", track.RegisteredAt},
									{"updated_at", now.UTC()},
								},
							)

							doc, err := m.outCollection.UpdateOne(
								context.Background(), filter, update, m.mongoService.GetUpsertOptions())

							ctx.JSON(200, Gin.H{
								"request":    track,
								"updated_at": now.UTC(),
								"stats": Gin.H{
									"new_id": doc.UpsertedID,
									"modified": doc.ModifiedCount,
									"upserted": doc.UpsertedCount,
									"matched": doc.MatchedCount,
								},
							})
						},
					},
				},
			},
		)
}
