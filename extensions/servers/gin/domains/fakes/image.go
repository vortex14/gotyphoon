package fakes

import (
	"bytes"
	"context"
	"fmt"
	Fake "github.com/vortex14/gotyphoon/extensions/data/fake"
	"image"
	"image/color"
	"net/http"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/fogleman/gg"
	Gin "github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
	"github.com/golang/freetype/truetype"

	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/interfaces"

	"github.com/vortex14/gotyphoon/elements/models/task"
	netHttp "github.com/vortex14/gotyphoon/extensions/pipelines/http/net-http"
	FakeImage "github.com/vortex14/gotyphoon/extensions/pipelines/image/fake-image"
	"github.com/vortex14/gotyphoon/extensions/servers/gin"
	GinExtension "github.com/vortex14/gotyphoon/extensions/servers/gin"
)


func CreateImageAction() interfaces.ActionInterface {
	return &GinExtension.Action{
		Action: &forms.Action{
			MetaInfo: &label.MetaInfo{
				Name:        "Fake image controller",
				Path:        FakeImagePath,
				Description: "Fake Image",
			},
			Methods:     []string{interfaces.GET},

			Pipeline: &forms.PipelineGroup{
				MetaInfo: &label.MetaInfo{
					Name:        "Fake image pipeline group",
					Description: "Creating fake image on request",
				},
				Stages: []interfaces.BasePipelineInterface{
					&gin.RequestPipeline{
						BasePipeline: &forms.BasePipeline{
							MetaInfo: &label.MetaInfo{
								Name:        "Request",
								Description: "Handle new request",
							},
							Middlewares: []interfaces.MiddlewareInterface{

},
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
					netHttp.CreatePrepareRequestPipeline(),
					netHttp.CreateRequestPipeline(),
					&netHttp.HttpResponsePipeline{
						BasePipeline: &forms.BasePipeline{
							MetaInfo: &label.MetaInfo{
								Name: "Response pipeline",
							},
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
							MetaInfo: &label.MetaInfo{
								Name:        "Prepare watermark",
								Description: "create rectangle on response image for watermark",
							},
						},
						Fn: func(context context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface, imgCtx *gg.Context) (error, context.Context) {
							logger.Info("Create watermark")
							w := float64(imgCtx.Width())
							h := float64(imgCtx.Height () /4)

							imgCtx.SetColor(color.RGBA{0, 0, 0, 204})
							imgCtx.DrawRectangle(0, float64(imgCtx.Height()-200), w, h)


							return nil, context
						},
					},
					&FakeImage.ImagePipeline{
						BasePipeline: &forms.BasePipeline{
							MetaInfo: &label.MetaInfo{
								Name:        "Create text on image",
								Description: "Create text in field on image",
							},
						},
						Fn: func(context context.Context, task interfaces.TaskInterface, logger interfaces.LoggerInterface, imgCtx *gg.Context) (error, context.Context) {

							box := packr.NewBox(".")
							source, _ := box.FindString("OpenSans-Bold.ttf")

							f, err := truetype.Parse([]byte(source))
							if err != nil {
								return err, nil
							}
							face := truetype.NewFace(f, &truetype.Options{
								Size: 90,
							})

							imgCtx.SetFontFace(face)

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
							y := float64(imgCtx.Height ( ) -70)
							imgCtx.DrawString(WATERMARK, x, y)
							imgCtx.Fill()

							return nil, context
						},
					},
					&gin.RequestPipeline{
						BasePipeline: &forms.BasePipeline{
							MetaInfo: &label.MetaInfo{
								Name:        "Result of Image",
								Description: "Send created image to client",
							},
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
		},
	}
}
