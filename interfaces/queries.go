package interfaces

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type MongoQuery struct {
	Query bson.D
	Group string
	Filter bson.D
	Database string
	Collection string
	Options interface{}
	Timeout time.Duration
	Context context.Context
}
