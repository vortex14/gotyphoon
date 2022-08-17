package forms

import (
	"context"
	Errors "errors"
	"testing"

	"github.com/vortex14/gotyphoon/elements/models/label"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

func init() {
	log.InitD()
}

func TestSkipStages(t *testing.T) {
	Convey("skip stages", t, func() {
		pg := &PipelineGroup{
			MetaInfo: &label.MetaInfo{
				Name:     "Skip-Pipeline",
				Required: true,
			},
			Stages: []interfaces.BasePipelineInterface{
				&BasePipeline{
					MetaInfo: &label.MetaInfo{
						Name:     "stage-1",
						Required: true,
					},
					Fn: func(ctx context.Context, logger interfaces.LoggerInterface) (error, context.Context) {
						logger.Info("run stage 1")
						ctx = PatchCtxPipelineGOTO(ctx, 4)
						return nil, ctx
					},
				},
				&BasePipeline{
					MetaInfo: &label.MetaInfo{
						Name:     "stage-2",
						Required: true,
					},
					Fn: func(ctx context.Context, logger interfaces.LoggerInterface) (error, context.Context) {
						logger.Info("run stage 2")
						return nil, ctx
					},
				},
				&BasePipeline{
					MetaInfo: &label.MetaInfo{
						Name:     "stage-3",
						Required: true,
					},
					Fn: func(ctx context.Context, logger interfaces.LoggerInterface) (error, context.Context) {
						logger.Info("run stage 3")
						return nil, ctx
					},
				},
				&BasePipeline{
					MetaInfo: &label.MetaInfo{
						Name:     "stage-4",
						Required: true,
					},
					Fn: func(ctx context.Context, logger interfaces.LoggerInterface) (error, context.Context) {
						logger.Info("run stage 4!")
						return nil, ctx
					},
				},
			},
		}

		err := pg.Run(context.Background())
		So(err, ShouldBeNil)
	})

}

func TestPanicPipeline(t *testing.T) {
	l := log.New(map[string]interface{}{"test": "test"})
	ctx := log.NewCtx(context.Background(), l)
	Convey("Create a pipeline with panic operation", t, func() {
		Pipe := &BasePipeline{
			MetaInfo: nil,
			Fn: func(ctx context.Context, logger interfaces.LoggerInterface) (error, context.Context) {
				logger.Debug("Run")

				var test map[string]interface{}

				test["test_panic_write_value"] = true

				return nil, nil
			},
			Cn: func(ctx context.Context, logger interfaces.LoggerInterface, err error) {
				logger.Error(err)
			},
		}
		var err error
		Pipe.Run(ctx, func(pipeline interfaces.BasePipelineInterface, error error) {
			err = error

		}, func(ctx context.Context) {
		})
		l.Debug(err)
		So(err, ShouldBeError)
	})
}

func TestRetryPipeline(t *testing.T) {
	l := log.New(map[string]interface{}{"test": "test"})
	ctx := log.NewCtx(context.Background(), l)

	Convey("Create a pipeline with error operation", t, func() {
		countIter := 0
		Pipe := &BasePipeline{
			Options:  GetDefaultRetryOptions(),
			MetaInfo: nil,
			Fn: func(ctx context.Context, logger interfaces.LoggerInterface) (error, context.Context) {
				logger.Debug("Run")
				countIter += 1
				return Errors.New("error operation"), nil
			},
			Cn: func(ctx context.Context, logger interfaces.LoggerInterface, err error) {
				logger.Error(err)
			},
		}
		var err error
		Pipe.Run(ctx, func(pipeline interfaces.BasePipelineInterface, error error) {
			err = error

		}, func(ctx context.Context) {
		})
		l.Debug(err)
		So(countIter, ShouldEqual, 7)
		So(err, ShouldBeError)
	})

	Convey("Create a pipeline with error operation and only once retry", t, func() {
		countIter := 0
		Pipe := &BasePipeline{
			Options:  GetNotRetribleOptions(),
			MetaInfo: nil,
			Fn: func(ctx context.Context, logger interfaces.LoggerInterface) (error, context.Context) {
				logger.Debug("Run")
				countIter += 1
				return Errors.New("error operation"), nil
			},
			Cn: func(ctx context.Context, logger interfaces.LoggerInterface, err error) {
				logger.Error(err)
			},
		}
		var err error
		Pipe.Run(ctx, func(pipeline interfaces.BasePipelineInterface, error error) {
			err = error

		}, func(ctx context.Context) {
		})
		l.Debug(err)
		So(countIter, ShouldEqual, 1)
		So(err, ShouldBeError)
	})

}
