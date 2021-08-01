package fake

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/vortex14/gotyphoon/data"
	"github.com/vortex14/gotyphoon/utils"
	"math/rand"
)

func init()  {
	level := 1
	gofakeit.AddFuncLookup("categories", gofakeit.Info{
		Category:    fmt.Sprintf("custom category level - %d", level),
		Description: "Random set categories",
		Output:      "item",

		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {

			u := utils.Utils{}
			var categoryName string

			switch level {
			case 1:
				categoryName = u.GetRandomFromSlice(data.CategoriesFirstLevel)
			case 2:
				categoryName = u.GetRandomFromSlice(data.CategoriesSecondLevel)
			case 3:
				categoryName = u.GetRandomFromSlice(data.CategoriesThirdLevel)
			}

			level += 1
			return categoryName, nil
		},
	})

	gofakeit.AddFuncLookup("categories_ids", gofakeit.Info{
		Category:    fmt.Sprintf("custom category level - %d", level),
		Description: "Random set images",
		Output:      "item",

		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			u := utils.Utils{}
			return u.GetRandomString(5, "012345"), nil
		},
	})
}
