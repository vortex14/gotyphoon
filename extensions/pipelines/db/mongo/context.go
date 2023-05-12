package mongo

import (
	"context"
	"fmt"
	"github.com/vortex14/gotyphoon/ctx"
	"github.com/vortex14/gotyphoon/integrations/mongo"

	M "go.mongodb.org/mongo-driver/mongo"
)

const (
	SERVICEMongo = "service_mongo"
	DATABASE     = "database"
	COLLECTION   = "collection"
)

func GetCollection(name string, context context.Context) (bool, *M.Collection) {
	collection, ok := ctx.Get(context, COLLECTION).(*M.Collection)
	return ok, collection

}

func SetCollection(context context.Context, collection *M.Collection) context.Context {
	return ctx.Update(context, COLLECTION, collection)
}

func GetDatabase(name string, context context.Context) (bool, *M.Database) {
	db, ok := ctx.Get(context, fmt.Sprintf("%s_%s", name, DATABASE)).(*M.Database)
	return ok, db
}

func SetDatabase(context context.Context, name string, db *M.Database) context.Context {
	return ctx.Update(context, fmt.Sprintf("%s_%s", name, DATABASE), db)
}

func GetService(context *context.Context) (bool, *mongo.Service) {
	db, ok := ctx.Get(*context, SERVICEMongo).(*mongo.Service)
	return ok, db
}

func SetService(context *context.Context, service *mongo.Service) context.Context {
	return ctx.Update(*context, SERVICEMongo, service)
}
