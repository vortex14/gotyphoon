package fake

import (
	"github.com/fatih/color"
	"math/rand"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/vortex14/gotyphoon/utils"
)

type Customer struct {
	Id string `fake:"{customer_id}" json:"id"`
	Name string `fake:"{name}" json:"name"`
	Phone string `fake:"+{phone}" json:"phone"`
	Email string `fake:"{email}" json:"email"`
}

func CreateCustomer() *Customer {
	var c *Customer
	err := gofakeit.Struct(&c)
	if utils.NotNill(err) {
		color.Red("%s", err.Error())
		return nil
	}
	return c
}


func init() {
	gofakeit.AddFuncLookup("customer_id", gofakeit.Info{
		Category:    "custom",
		Description: "Random set customer_id",
		Output:      "str",
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			u := utils.Utils{}
			return u.GetUUID(), nil
		},
	})
}