package mongo

import (
	"github.com/vortex14/gotyphoon/integrations/mongo"
	"github.com/vortex14/gotyphoon/interfaces"

	M "go.mongodb.org/mongo-driver/mongo"

	"context"
)

type Callback func(
	context context.Context,
	task interfaces.TaskInterface,
	logger interfaces.LoggerInterface,
	service *mongo.Service,
	database *M.Database,
	collection *M.Collection) (error, context.Context)

type CancelCallback func(
	err error,
	context context.Context,
	task interfaces.TaskInterface,
	logger interfaces.LoggerInterface,
)
