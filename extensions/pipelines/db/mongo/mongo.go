package mongo

import (
	Context "context"
	"fmt"
	"sync"

	M "go.mongodb.org/mongo-driver/mongo"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/task"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/extensions/pipelines"
	"github.com/vortex14/gotyphoon/integrations/mongo"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

type Pipeline struct {
	*forms.BasePipeline
	*pipelines.TaskPipeline

	opts *interfaces.ServiceMongo

	ctx  Context.Context
	sCtx sync.Once

	Fn func(
		context Context.Context,
		task interfaces.TaskInterface,
		logger interfaces.LoggerInterface,

		service *mongo.Service,
		database *M.Database,
		collection *M.Collection,

	) (error, Context.Context)

	Cn func(
		err error,
		context Context.Context,
		task interfaces.TaskInterface,
		logger interfaces.LoggerInterface,
	)
}

func (t *Pipeline) UnpackResponseCtx(
	ctx Context.Context,
) (bool, interfaces.TaskInterface, interfaces.LoggerInterface, *mongo.Service, *M.Database, *M.Collection) {
	okT, taskInstance := task.Get(ctx)
	okL, logger := log.Get(ctx)

	var okS bool
	var db *M.Database
	var cl *M.Collection
	var sv *mongo.Service

	if t.SharedCtxStatus {
		okS, sv = GetService(t.SharedCtx)
	} else {
		okS, sv = GetService(&t.ctx)
	}

	if sv != nil {
		_, db = GetDatabase(t.opts.DefaultDatabase, t.ctx)

		_, cl = GetCollection(t.opts.DefaultCollection, t.ctx)

	}

	return okL && okT && okS, taskInstance, logger, sv, db, cl
}

func (t *Pipeline) Run(
	context Context.Context,
	reject func(pipeline interfaces.BasePipelineInterface, err error),
	next func(ctx Context.Context),
) {

	if t.Fn == nil {
		reject(t, Errors.TaskPipelineRequiredHandler)
		return
	}

	if t.opts != nil {
		t.initCtx()
	}

	ok, taskInstance, logger, sv, db, collection := t.UnpackResponseCtx(context)

	if !ok {

		fError := fmt.Errorf("%s. taskInstance: %s, logger: %s, db: %+v, collection: %v, sv: %+v",
			Errors.PipelineContexFailed, taskInstance, logger, db, collection, sv)
		reject(t, fError)
		t.Cancel(context, logger, fError)
		return
	}

	t.SafeRun(context, logger, func(patchedCtx Context.Context) error {

		err, newContext := t.Fn(patchedCtx, taskInstance, logger, sv, db, collection)
		if err != nil {
			return err
		}
		next(newContext)
		return nil

	}, func(err error) {

		reject(t, err)

	})

}

func (t *Pipeline) initCtx() {
	t.sCtx.Do(func() {
		t.ctx = Context.Background()

		var service *mongo.Service

		if t.SharedCtxStatus {
			s, srv := GetService(t.SharedCtx)
			if s {
				service = srv
			} else {
				service = &mongo.Service{
					Settings: *t.opts,
				}
				*t.SharedCtx = SetService(t.SharedCtx, service)
			}
		} else {
			service = &mongo.Service{
				Settings: *t.opts,
			}
			t.ctx = SetService(&t.ctx, service)
		}

		if len(t.opts.DefaultCollection) > 0 && len(t.opts.DefaultDatabase) > 0 {
			t.ctx = SetDatabase(t.ctx, t.opts.DefaultDatabase, service.GetMongoDB(t.opts.DefaultDatabase))
			t.ctx = SetCollection(t.ctx, service.GetMongoCollection(t.opts.DefaultDatabase, t.opts.DefaultCollection))
		}

	})
}

func (t *Pipeline) Cancel(
	context Context.Context,
	logger interfaces.LoggerInterface,
	err error,
) {

	if t.Cn == nil {
		return
	}

	ok, taskInstance, logger := t.UnpackCtx(context)
	if !ok {
		return
	}

	t.Cn(err, context, taskInstance, logger)

}
