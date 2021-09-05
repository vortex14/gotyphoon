package fake

import (
	"fmt"
	"github.com/vortex14/gotyphoon/extensions/data"
	"math/rand"

	"github.com/fatih/color"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/vortex14/gotyphoon/utils"
)

type Category struct {
	Name string `fake:"{categories}" json:"name"`
	Id   string `fake:"{categories_ids}" json:"id"`
}

type Categories struct {
	Categories []string `fake:"{categories}" fakesize:"3" json:"categories"`
	CategoriesIds []string `fake:"{categories_ids}" fakesize:"3" json:"categories_ids"`
}

func CreateCategories() *Categories {
	var c *Categories
	err := gofakeit.Struct(&c)
	if utils.NotNill(err) {
		color.Red("%s", err.Error())
		return nil
	}
	return c
}

func CreateCategory() *Category {
	var c *Category
	err := gofakeit.Struct(&c)
	if utils.NotNill(err) {
		color.Red("%s", err.Error())
		return nil
	}
	return c
}


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
