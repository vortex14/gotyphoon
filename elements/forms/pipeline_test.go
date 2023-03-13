package forms

import (
	"github.com/avast/retry-go/v4"
	"github.com/google/uuid"
	"golang.org/x/sync/semaphore"

	"context"
	Errors "errors"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/vortex14/gotyphoon/elements/models/label"
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
			l.Error(error)
			err = error

		}, func(ctx context.Context) {
		})
		l.Debug(err)
		So(err, ShouldBeError)
	})
}

func TestPipelineRetry(t *testing.T) {
	l := log.New(map[string]interface{}{"test": "test"})
	ctx := log.NewCtx(context.Background(), l)
	Convey("Create a pipeline with error operation and check retry process", t, func() {
		countIter := 0
		Pipe := &BasePipeline{
			Options:  &Options{Retry: RetryOptions{MaxCount: 7}, MaxConcurrent: 2},
			MetaInfo: &label.MetaInfo{Name: "base-pipeline"},
			Fn: func(ctx context.Context, logger interfaces.LoggerInterface) (error, context.Context) {
				logger.Debug("Run")
				countIter += 1
				return Errors.New("error operation"), nil
			},
			Cn: func(ctx context.Context, logger interfaces.LoggerInterface, err error) {
				logger.Error(err)
				So(err, ShouldBeError)
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
}

func TestRetry(t *testing.T) {

	Convey("retry", t, func() {
		count := 0

		testCallback := func() error {
			count += 1
			return Errors.New("test main error")
		}
		e := retry.Do(testCallback, retry.Attempts(uint(10)))

		So(e, ShouldBeError)

		So(count, ShouldEqual, 10)

	})

}

func TestRetryDelayPipeline(t *testing.T) {
	l := log.New(map[string]interface{}{"test": "test"})
	ctx := log.NewCtx(context.Background(), l)
	Convey("Create a pipeline with error operation and check retry process with delay", t, func() {
		countIter := 0
		Pipe := &BasePipeline{
			Options:  &Options{Retry: RetryOptions{MaxCount: 7, Delay: time.Second * 3}, MaxConcurrent: 1},
			MetaInfo: &label.MetaInfo{Name: "base-pipeline"},
			Fn: func(ctx context.Context, logger interfaces.LoggerInterface) (error, context.Context) {
				logger.Debug("Run")
				countIter += 1
				return Errors.New("error operation"), nil
			},
			Cn: func(ctx context.Context, logger interfaces.LoggerInterface, err error) {
				logger.Error(err)
				So(err, ShouldBeError)
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
}

func TestSemPipeline(t *testing.T) {
	//cuncurrentCall := 0
	p := BasePipeline{
		Options: &Options{MaxConcurrent: 1},
		Fn: func(ctx context.Context, logger interfaces.LoggerInterface) (error, context.Context) {
			logger.Info("Run")
			time.Sleep(10 * time.Second)

			return nil, ctx
		},
	}

	for i := 0; i < 10; i++ {
		l := log.New(map[string]interface{}{"number": i, "uuid": uuid.New().String()})
		ctx := log.NewCtx(context.Background(), l)
		go p.Run(ctx, func(pipeline interfaces.BasePipelineInterface, err error) {
			l.Error(err)
		}, func(ctx context.Context) {
			l.Debug("a good data")
		})

		//if !p.Try() {
		//	l.Error("Busy !!")
		//}

	}

	time.Sleep(60 * time.Second)

	//p.Await()

}

func TestSem(t *testing.T) {
	s := semaphore.NewWeighted(10)

	e := s.Acquire(context.Background(), 1)

	println(e)
}
