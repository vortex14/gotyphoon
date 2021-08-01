package fake

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/vortex14/gotyphoon/data"
	"github.com/vortex14/gotyphoon/utils"
	"math/rand"
)


type ProductShippingDimensions struct {
	Dimensions struct {
		Height struct {
			Value float64 `json:"value" fake:"{value_size}"`
			Unit  string `json:"unit" fake:"{unit}"`
		} `json:"height"`
		Length struct {
			Value float64 `json:"value" fake:"{value_size}"`
			Unit  string `json:"unit" fake:"{unit}"`
		} `json:"length"`
		Depth struct {
			Value float64 `json:"value" fake:"{value_size}"`
			Unit  string `json:"unit" fake:"{unit}"`
		} `json:"depth"`
		Weight struct {
			Value float64 `json:"value" fake:"{value_size}"`
			Unit  string `json:"unit" fake:"{unit_w}"`
		} `json:"weight"`
	} `json:"dimensions"`
}


type Product struct {
	Price struct {
		OfferPrice float32 `fake:"{price:0,100}" json:"offerPrice"`
		ListingPrice float32 `fake:"{price:100,200}" json:"listingPrice"`
	} `json:"price"`
	Url string `fake:"{url}" json:"url"`
	Upc string `fake:"{upc}" json:"upc"`
	Color string `fake:"{color}" json:"color"`
	Brand string `fake:"{brand}" json:"brand"`
	Title string `fake:"{sentence}" json:"title"`
	Quantity int `fake:"{number:1,10}" json:"quantity"`
	Shipping ProductShippingDimensions `json:"shipping"`
	ProductId string `fake:"{product_id}" json:"productId"`
	Description string `fake:"{paragraph}" json:"description"`
	Images []string `fake:"{images}" fakesize:"3" json:"images"`
	Categories []string `fake:"{categories}" fakesize:"3" json:"categories"`
	ApiProvider string `json:"api_provider" fake:"{randomstring:[typhoon]}"`
	CountryOfOrigin string `fake:"{randomstring:[USA,CA]}" json:"countryoforigin"`
	CategoriesIds []string `fake:"{categories_ids}" fakesize:"3" json:"categories_ids"`
	Marketplace string `fake:"{randomstring:[ebay.com,amazon.com,walmart.com,homedepot.com]}" json:"marketplace"`
}



func init()  {
	gofakeit.AddFuncLookup("product_id", gofakeit.Info{
		Category:    "custom",
		Description: "Random set images",
		Output:      "list",
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			u := utils.Utils{}
			return u.GetUUID(), nil
		},
	})


	gofakeit.AddFuncLookup("value_size", gofakeit.Info{
		Category:    "custom",
		Description: "Random set images",
		Output:      "list",
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {

			u := utils.Utils{}

			return u.GetRandomFloat(), nil
		},
	})
	gofakeit.AddFuncLookup("unit", gofakeit.Info{
		Category:    "custom",
		Description: "Random set images",
		Output:      "list",
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			u := utils.Utils{}
			return u.GetRandomFromSlice([]string{"m"}), nil
		},
	})

	gofakeit.AddFuncLookup("unit_w", gofakeit.Info{
		Category:    "custom",
		Description: "Random set images",
		Output:      "list",
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			u := utils.Utils{}
			return u.GetRandomFromSlice([]string{"g"}), nil
		},
	})

	gofakeit.AddFuncLookup("brand", gofakeit.Info{
		Category:    "custom",
		Description: "Random set images",
		Output:      "list",
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			u := utils.Utils{}
			brand := u.GetRandomFromSlice(data.Brands)
			return brand, nil
		},
	})
}
