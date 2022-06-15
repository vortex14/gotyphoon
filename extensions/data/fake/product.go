package fake

import (
	"github.com/vortex14/gotyphoon/extensions/data"
	"math/rand"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/fatih/color"

	"github.com/vortex14/gotyphoon/utils"
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

type BaseProduct struct {
	Url string `fake:"{url}" json:"url"`
	Upc string `fake:"{upc}" json:"upc"`
	Id string `fake:"{product_id}" json:"id"`

	Title string `fake:"{sentence}" json:"title"`
	Description string `fake:"{paragraph}" json:"description"`
	Price struct {
		OfferPrice float32 `fake:"{price:0,100}" json:"offerPrice"`
		ListingPrice float32 `fake:"{price:100,200}" json:"listingPrice"`
	} `json:"price"`
}

type StockProduct struct {
	Quantity  int                       `fake:"{number:1,10}" json:"quantity"`
}

type MediaProductDetails struct {
	Images []string `fake:"{images}" fakesize:"3" json:"images"`
}

type ProviderDetails struct {
	ApiProvider string `json:"api_provider" fake:"{randomstring:[typhoon]}"`
	CountryOfOrigin string `fake:"{randomstring:[USA,CA]}" json:"countryoforigin"`
	Marketplace string `fake:"{randomstring:[ebay.com,amazon.com,walmart.com,homedepot.com]}" json:"marketplace"`
}

type ProductAttributes struct {
	Color string `fake:"{color}" json:"color"`
	Brand string `fake:"{brand}" json:"brand"`
}

type Product struct {
	BaseProduct
	Categories
	StockProduct
	ProviderDetails



	
	Shipping  ProductShippingDimensions `json:"shipping"`
	ProductId string                    `fake:"{product_id}" json:"productId"`
	
	
}

func CreateProduct() *Product {
	var p *Product
	err := gofakeit.Struct(&p)
	if utils.NotNill(err) {
		color.Red("%s", err.Error())
		return nil
	}
	return p
}

func CreateProductWithId() *BaseProduct {
	var p *BaseProduct
	err := gofakeit.Struct(&p)
	if utils.NotNill(err) {
		color.Red("%s", err.Error())
		return nil
	}
	return p
}


func CreateShipping() *ProductShippingDimensions {
	var p *ProductShippingDimensions
	err := gofakeit.Struct(&p)
	if utils.NotNill(err) {
		color.Red("%s", err.Error())
		return nil
	}
	return p
}


func init()  {
	gofakeit.AddFuncLookup("product_id", gofakeit.Info{
		Category:    "custom",
		Description: "Random set product_id",
		Output:      "list",
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			u := utils.Utils{}
			return u.GetUUID(), nil
		},
	})


	gofakeit.AddFuncLookup("value_size", gofakeit.Info{
		Category:    "custom",
		Description: "Random set value_size",
		Output:      "list",
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {

			u := utils.Utils{}

			return u.GetRandomFloat(), nil
		},
	})
	gofakeit.AddFuncLookup("unit", gofakeit.Info{
		Category:    "custom",
		Description: "Random set unit",
		Output:      "list",
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			u := utils.Utils{}
			return u.GetRandomFromSlice([]string{"m"}), nil
		},
	})

	gofakeit.AddFuncLookup("unit_w", gofakeit.Info{
		Category:    "custom",
		Description: "Random set weight",
		Output:      "list",
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			u := utils.Utils{}
			return u.GetRandomFromSlice([]string{"g"}), nil
		},
	})

	gofakeit.AddFuncLookup("brand", gofakeit.Info{
		Category:    "custom",
		Description: "Random set brand",
		Output:      "list",
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			u := utils.Utils{}
			brand := u.GetRandomFromSlice(data.Brands)
			return brand, nil
		},
	})
}
