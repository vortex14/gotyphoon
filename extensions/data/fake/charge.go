package fake

import (
	"github.com/fatih/color"
	"math/rand"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/vortex14/gotyphoon/utils"
)

type Charge struct {
	Id string `fake:"{charge_id}" json:"id"`
	Amount float32 `fake:"{price}" json:"amount"`
	Currency string `fake:"{currency}" json:"currency"`
}


func CreateCharge() *Charge {
	var c *Charge
	err := gofakeit.Struct(&c)
	if utils.NotNill(err) {
		color.Red("%s", err.Error())
		return nil
	}
	return c
}


func init() {
	gofakeit.AddFuncLookup("charge_id", gofakeit.Info{
		Category:    "custom",
		Description: "Random set charge_id",
		Output:      "str",
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			u := utils.Utils{}
			return u.GetUUID(), nil
		},
	})
}