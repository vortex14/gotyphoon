package gif

import (
	"fmt"
	"github.com/vortex14/gotyphoon/log"
	"image"

	"github.com/vortex14/gotyphoon/elements/models/awaitabler"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/elements/models/singleton"
	"github.com/vortex14/gotyphoon/utils/img"
)

type Gif struct {
	label.MetaInfo
	awaitabler.Object
	singleton.Singleton

	Fps          int
	PathImg      []string
	ExportPath   string
	SourceFormat string
	images       []image.Image
}

func (g *Gif) Append(path string) *Gif {
	g.PathImg = append(g.PathImg, path)
	return g
}

func (g *Gif) Create() {
	g.Construct(func() {
		logger := log.New(log.D{"gif": "builder"})
		err := img.BuildGif(g.PathImg, g.Fps, g.ExportPath)
		if err != nil {
			logger.Error(err.Error())
		} else {
			logger.Debug(fmt.Sprintf("Gif %s created !", g.ExportPath))
		}
	})
}
