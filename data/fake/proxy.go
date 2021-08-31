package fake

import (
	"fmt"
	browser "github.com/EDDYCJY/fake-useragent"
	"math/rand"
	"reflect"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/bxcodec/faker"
	"github.com/fatih/color"
)

type Proxy struct {
	Proxy string `fake:"{proxy}" json:"proxy"`
	Agent string `fake:"{useragent}" json:"agent"`
	Success bool `fake:"{success}" json:"success"`
	AgentMobile string `fake:"{mobile}" json:"agent_mobile"`
}

func init()  {
	gofakeit.AddFuncLookup("proxy", gofakeit.Info{
		Category:    "custom",
		Description: "Random set proxy",
		Output:      "str",
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			inet := faker.Internet{}

			data, err := inet.IPv4(reflect.Value{})
			if err != nil {
				return nil, err
			}

			ip := fmt.Sprintf("https://%s:3128", data.(string))

			return ip, nil
		},
	})

	gofakeit.AddFuncLookup("success", gofakeit.Info{
		Category:    "custom",
		Description: "Status proxy server",
		Output:      "bool",
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			return true, nil
		},
	})

	gofakeit.AddFuncLookup("useragent", gofakeit.Info{
		Category:    "custom",
		Description: "set random useragent",
		Output:      "str",
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			return browser.Random(), nil
		},
	})

	gofakeit.AddFuncLookup("mobile", gofakeit.Info{
		Category:    "custom",
		Description: "set random mobile agent",
		Output:      "str",
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			return browser.Mobile(), nil
		},
	})
}

func CreateFakeProxy() *Proxy {
	var proxy *Proxy
	err := gofakeit.Struct(&proxy)
	if err != nil {
		color.Red("%s", err.Error())
		return nil
	}
	return proxy
}