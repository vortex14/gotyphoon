package forms

import (
	"github.com/avast/retry-go/v4"
	"github.com/google/uuid"
	"github.com/vortex14/gotyphoon/log"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"

	"context"
	Errors "errors"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/interfaces"
)

func TestNewPipeline(t *testing.T) {
	pl := &BasePipeline{MetaInfo: &label.MetaInfo{Name: "stage-1", Required: true}, Fn: func(ctx context.Context, logger interfaces.LoggerInterface) (error, context.Context) {
		logger.Info("run stage 1")
		return nil, ctx
	}}
	if pl.Name != "stage-1" {
		t.Fatal("name isn't valid")
	}
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

func TestPipelineLabel(t *testing.T) {
	l := log.New(log.DebugLevel, map[string]interface{}{"test": "test"})
	ctx := log.NewCtx(context.Background(), l)
	Convey("test pipeline label", t, func() {

		Pipe := &BasePipeline{
			MetaInfo: &label.MetaInfo{Name: "test-1-3-4"},
			Fn: func(ctx context.Context, logger interfaces.LoggerInterface) (error, context.Context) {
				logger.Debug("Run")

				_, _labels := GetPipelineLabel(ctx)

				So(_labels.Name, ShouldEqual, "test-1-3-4")

				return nil, nil
			},
			Cn: func(ctx context.Context, logger interfaces.LoggerInterface, err error) {
				logger.Error("pipeline", zap.Error(err))
			},
		}
		var err error
		Pipe.Run(ctx, func(context context.Context, pipeline interfaces.BasePipelineInterface, error error) {

		}, func(ctx context.Context) {
		})

		So(err, ShouldBeNil)
	})

}

func TestPanicPipeline(t *testing.T) {
	l := log.New(log.DebugLevel, map[string]interface{}{"test": "test"})
	ctx := log.NewCtx(context.Background(), l)
	Convey("Create a pipeline with panic operation", t, func() {
		var errCallback error
		Pipe := &BasePipeline{
			MetaInfo: &label.MetaInfo{Name: "panic pipe"},
			Fn: func(ctx context.Context, logger interfaces.LoggerInterface) (error, context.Context) {
				logger.Debug("Run")

				var test map[string]interface{}

				test["test_panic_write_value"] = true

				return nil, nil
			},
			Cn: func(ctx context.Context, logger interfaces.LoggerInterface, err error) {
				errCallback = err
				logger.Error("pipeline.cn", zap.Error(err))
			},
		}
		var err error
		Pipe.Run(ctx, func(context context.Context, pipeline interfaces.BasePipelineInterface, error error) {
			l.Error("err", zap.Error(error))
			err = error

		}, func(ctx context.Context) {
		})

		So(err, ShouldBeError)
		So(errCallback, ShouldBeError)
	})
}

func TestPipelineRetry(t *testing.T) {
	l := log.New(log.DebugLevel, map[string]interface{}{"test": "test"})
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
				logger.Error("pipeline.cn", zap.Error(err))
				So(err, ShouldBeError)
			},
		}
		var err error
		Pipe.Run(ctx, func(context context.Context, pipeline interfaces.BasePipelineInterface, error error) {
			err = error

		}, func(ctx context.Context) {
		})
		l.Debug("err", zap.Error(err))
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
	l := log.New(log.DebugLevel, map[string]interface{}{"test": "test"})
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
				logger.Error("pipeline.cn", zap.Error(err))
				So(err, ShouldBeError)
			},
		}
		var err error
		Pipe.Run(ctx, func(context context.Context, pipeline interfaces.BasePipelineInterface, error error) {
			err = error

		}, func(ctx context.Context) {
		})
		l.Debug("error", zap.Error(err))
		So(countIter, ShouldEqual, 7)
		So(err, ShouldBeError)
	})
}

func TestSemPipeline(t *testing.T) {

	Convey("test semaphore group", t, func() {
		p := BasePipeline{
			MetaInfo: &label.MetaInfo{Name: "semaphore pipe"},
			Options:  &Options{MaxConcurrent: 1},
			Fn: func(ctx context.Context, logger interfaces.LoggerInterface) (error, context.Context) {
				logger.Info("Run")
				time.Sleep(10 * time.Second)

				return nil, ctx
			},
		}

		crowded := 0
		success := 0

		for i := 0; i < 10; i++ {

			l := log.New(log.DebugLevel, map[string]interface{}{"number": i, "uuid": uuid.New().String()})
			ctx := log.NewCtx(context.Background(), l)
			go p.Run(ctx, func(context context.Context, pipeline interfaces.BasePipelineInterface, err error) {
				crowded++
				l.Error("p.err", zap.Error(err))
			}, func(ctx context.Context) {
				success++
				l.Debug("a good data")
			})
		}

		time.Sleep(11 * time.Second)

		So(success, ShouldEqual, 1)
		So(crowded, ShouldEqual, 9)

	})

}

func TestSem(t *testing.T) {
	s := semaphore.NewWeighted(10)

	e := s.Acquire(context.Background(), 1)

	println(e)
}
