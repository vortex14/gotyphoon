package fake

import (
	"fmt"
	"math/rand"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/vortex14/gotyphoon/data"
	"github.com/vortex14/gotyphoon/utils"

)

func init()  {
	gofakeit.AddFuncLookup("categories", gofakeit.Info{
		Category:    fmt.Sprintf("custom category level"),
		Description: "Random set categories",
		Output:      "item",

		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {

			u := utils.Utils{}

			level := utils.GetRandomIntRange(4,1)
			var categoryName string

			switch level {
			case 1:
				categoryName = u.GetRandomFromSlice(data.CategoriesFirstLevel)
			case 2:
				categoryName = u.GetRandomFromSlice(data.CategoriesSecondLevel)
			case 3:
				categoryName = u.GetRandomFromSlice(data.CategoriesThirdLevel)

			}
			return categoryName, nil
		},
	})

	gofakeit.AddFuncLookup("categories_ids", gofakeit.Info{
		Category:    fmt.Sprintf("custom category level"),
		Description: "Random set images",
		Output:      "item",

		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			u := utils.Utils{}
			return u.GetRandomString(5, "012345"), nil
		},
	})
}
