package fake

import (
	"math/rand"

	"github.com/brianvoe/gofakeit/v6"
)

func init()  {
	gofakeit.AddFuncLookup("images", gofakeit.Info{
		Category:    "custom",
		Description: "Random set images",
		Output:      "list",
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			const width = 1000
			const height = 500
			imageUrl := gofakeit.ImageURL(width, height)
			return imageUrl, nil
		},
	})
}

