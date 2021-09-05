package fake

import (
	"fmt"
	"github.com/vortex14/gotyphoon/utils"
	"math/rand"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/fatih/color"
	"github.com/osamingo/checkdigit"
)

var letters = []rune("0123456789")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}


func init()  {
	gofakeit.AddFuncLookup("upc", gofakeit.Info{
		Category:    "custom",
		Description: "Random UPC code",
		Example:     "bill",
		Output:      "string",
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {

			rand.Seed(time.Now().UnixNano())
			p := checkdigit.NewUPC()
			seed := randSeq(11)
			cd, err := p.Generate(seed)
			if utils.NotNill(err) {
				color.Red("failed to generate check digit")
				return nil, nil
			}

			ok := p.Verify(seed + strconv.Itoa(cd))
			if !ok {
				return nil, nil
			}
			upc := fmt.Sprintf("%s%d", seed, cd)
			//fmt.Printf("seed: %s, check digit: %d, verify: %t\n", seed, cd, ok)

			return upc, nil
		},
	})
}
