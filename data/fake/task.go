package fake

import (
	"encoding/base64"
	"math/rand"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/color"

	"github.com/vortex14/gotyphoon/interfaces"
	typhoonTask "github.com/vortex14/gotyphoon/task"
	"github.com/vortex14/gotyphoon/utils"
)

func init()  {
	gofakeit.AddFuncLookup("response_product", gofakeit.Info{
		Category:    "custom",
		Description: "Random set response product",
		Output:      "str",
		Generate: func(r *rand.Rand, m *gofakeit.MapParams, info *gofakeit.Info) (interface{}, error) {
			u := utils.Utils{}
			var f Product
			_ = gofakeit.Struct(&f)
			dump := u.PrintPrettyJson(f)
			encoded := base64.StdEncoding.EncodeToString([]byte(dump))
			return encoded, nil
		},
	})
}

func CreateDefaultTask() *typhoonTask.TyphoonTask {
	task, _ := CreateFakeTask(interfaces.FakeTaskOptions{})
	return task
}

func CreateFakeTask(options interfaces.FakeTaskOptions)  (*typhoonTask.TyphoonTask, error){
	// TODO: task.yaml
	var task typhoonTask.TyphoonTask
	err := gofakeit.Struct(&task)
	if err != nil {
		color.Red("%s", err.Error())
		return nil, err
	}

	if options.UserAgent {
		task.Fetcher.Headers = map[string]string{
			"User-Agent": browser.Chrome(),
		}
	}

	if !options.Proxy {
		task.Fetcher.Proxy = ""
	}

	if options.Cookies {
		jwtoken := jwt.New(jwt.SigningMethodHS256)

		jwtS, _ := jwtoken.SigningString()

		task.Fetcher.Cookies = map[string]string{
			"jwt": jwtS,
		}
	}

	if options.Auth {
		task.Fetcher.Auth = map[string]string{
			"login": gofakeit.Username(),
			"password": gofakeit.Password(true, false, true, false, false, 10),
		}
	}



	return &task, nil
}


