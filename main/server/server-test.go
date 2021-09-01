package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"net/http"
	"os"
	"path/filepath"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/fogleman/gg"
	Gin "github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	Fake "github.com/vortex14/gotyphoon/data/fake"
	"github.com/vortex14/gotyphoon/elements/forms"
	net_http "github.com/vortex14/gotyphoon/extensions/pipelines/http/net-http"
	FakeImage "github.com/vortex14/gotyphoon/extensions/pipelines/image/fake-image"
	"github.com/vortex14/gotyphoon/extensions/servers/pipelines/gin"
	"github.com/vortex14/gotyphoon/extensions/servers/resources/home"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/server"
	"github.com/vortex14/gotyphoon/task"
)

const (
	ImgCTX = "IMG_CTX"

	WATERMARK = "image.typhoon.dev"
)

func main()  {
	err := (&server.TyphoonServer{
		BaseServerLabel: &interfaces.BaseServerLabel{
			Name:        "Server image generator",
			Description: "Generator images",
		},
		Port: 17667,
		IsDebug: true,
	}).Init().InitLogger().AddResource(home.Constructor("/").AddAction(&interfaces.Action{
		Name:        "image",
		Description: "Image data faker",
		Path:        "image",
		Methods: []string{interfaces.GET},
		Pipeline: &forms.PipelineGroup{
			Stages: []interfaces.BasePipelineInterface{
				&gin.RequestPipeline{
					BasePipeline: &forms.BasePipeline{
						Name:        "Request",
						Description: "Handle new request",
					},
					Fn: func(context context.Context, ginCtx *Gin.Context, logger interfaces.LoggerInterface) (error, context.Context) {
						logger.Info("new request")
						imageUrl := gofakeit.ImageURL(1840, 1024)
						taskF := Fake.CreateDefaultTask()
						taskF.SetFetcherUrl(imageUrl)
						context = task.PatchCtx(context, taskF)
						return nil, context
					},
				},
				net_http.CreatePrepareRequestPipeline(),
				net_http.CreateRequestPipeline(),
				&net_http.HttpResponsePipeline{
					BasePipeline: &forms.BasePipeline{
						Name: "Response pipeline",
					},
					Fn: func(
						context context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface,
						client *http.Client, request *http.Request, transport *http.Transport,
						response *http.Response, data *string) (error, context.Context){

						r := bytes.NewReader([]byte(*data))
						img, _, _ := image.Decode(r)
						imgCtx := gg.NewContextForImage(img)
						newCtx := FakeImage.NewImgCtx(context, imgCtx)

						logger.Warning(fmt.Sprintf("response code: %d, len: %d", response.StatusCode, len(*data)))

						return nil, newCtx
					},
				},
				&FakeImage.ImagePipeline{
					BasePipeline: &forms.BasePipeline{
						Name:        "Prepare watermark",
						Description: "create rectangle on response image for watermark",
					},
					Fn: func(context context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface, imgCtx *gg.Context) (error, context.Context) {
						logger.Info("Create watermark")
						w := float64(imgCtx.Width())
						h := float64(imgCtx.Height()/4)

						imgCtx.SetColor(color.RGBA{0, 0, 0, 204})
						imgCtx.DrawRectangle(0, float64(imgCtx.Height()-200), w, h)


						return nil, context
					},
				},
				&FakeImage.ImagePipeline{
					BasePipeline: &forms.BasePipeline{
						Name:        "Create text on image",
						Description: "Create text in field on image",
					},
					Fn: func(context context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface, imgCtx *gg.Context) (error, context.Context) {

						root, _ := os.Getwd()
						fontPath := filepath.Join(root, "main", "server", "OpenSans-Bold.ttf")

						if err := imgCtx.LoadFontFace(fontPath, 90); err != nil {
							return err, nil
						}
						textColor := color.White

						_, g, b, _ := textColor.RGBA()

						mutedColor := color.RGBA{
							R: uint8(255),
							G: uint8(g),
							B: uint8(b),
							A: uint8(200),
						}

						imgCtx.SetColor(mutedColor)
						_, _ = imgCtx.MeasureString(WATERMARK)
						x := float64(70)
						y := float64(imgCtx.Height()-70)
						imgCtx.DrawString(WATERMARK, x, y)
						imgCtx.Fill()

						return nil, context
					},
				},
				&gin.RequestPipeline{
					BasePipeline: &forms.BasePipeline{
						Name:        "Result of Image",
						Description: "Send created image to client",
					},
					Fn: func(context context.Context, ginCtx *Gin.Context, logger interfaces.LoggerInterface) (error, context.Context) {
						logger.Info("send image")

						wr := &bytes.Buffer{}
						_, imgCtx := FakeImage.GetImgCtx(context)

						_ = imgCtx.EncodePNG(wr)

						_, _ = ginCtx.Writer.Write(wr.Bytes())


						return nil, context
					},
				},

			},	
		},
	})).Run()


	if err != nil {
		logrus.Error(err.Error())
	}
}
