package fake

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/fatih/color"

	"github.com/vortex14/gotyphoon/utils"
)

type Payment struct {
	Id string `fake:"{charge_id}" json:"id"`
	Price float32 `fake:"{price}" json:"price"`
	Currency string `fake:"{currency}" json:"currency"`
	CreditCardExp   string  `fake:"{creditcardexp}" json:"credit_card_exp"`
	CreditCardType  string `fake:"{creditcardtype}" json:"creadit_card_type"`
	CreditCardNumber string `fake:"{creditcardnumber}" json:"credit_card_number"`
}

func CreatePayment() *Payment {
	var c *Payment
	err := gofakeit.Struct(&c)
	if utils.NotNill(err) {
		color.Red("%s", err.Error())
		return nil
	}
	return c
}
