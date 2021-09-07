package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/vortex14/gotyphoon/extensions/data/fake"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
	"github.com/vortex14/gotyphoon/utils"
)

func init()  {
	log.InitD()
}


var productIds = []int{3454581, 3454582, 3470493, 3931128, 3628008, 3628009, 3687083, 4065768, 3640165, 3406204,
	3406206, 3406207, 3406208, 4095504, 4095502, 4095503, 3684244, 4086986, 3638745, 3662012, 3607978, 3607977,
	3639762, 3639763, 4088768, 4088767, 4102294, 3184015, 3184021, 3184013, 4039051, 4039052, 4039053, 3641346,
	3641347, 3651895, 3641348, 3651880, 3992299, 4060980, 4060981, 3472905, 3472906, 3469466, 3974462, 3671791,
	3675139, 3675180, 3673091, 3632746, 3632745, 4074575, 4074577, 4074576, 4074578, 3640382, 3659029, 3608102,
	3919341, 3401039, 3942977, 3671132, 3406685, 3966538, 3686188, 3639719, 3639740, 3969427, 4018759, 4015600,
	4025353, 3404339, 3659796, 3428589, 3385436, 3505113, 3680787, 3480408, 3305163, 4082115, 3604966, 3688990,
	3659421, 3668090, 3424514, 3685992, 3456174, 3427136, 2155855, 3445604, 4070251, 4070252, 4070253, 4070254,
	4086974, 4086976, 3194068, 3377483, 3377489, 3375417, 3375415, 3375416, 3429446, 3429447, 3677285, 3478733,
	3478734, 3658363, 3680932, 4005426, 4070213, 3971940, 4012920, 4058577, 3624885, 3624886, 3671814, 3624248,
	4072019, 4087024, 4072282, 4007658, 3641250, 3692211, 3694524, 3234684, 3911586, 3661417, 4087446, 4087447,
	4094340, 3450476, 3443522, 3445764, 3908200, 3908201, 3686124, 3686125, 2152649, 3369165, 4085999, 3630356,
	4072038, 3634639, 3405617, 3405618, 3698600, 69435, 75230, 75228, 3430225, 4075001, 4075005, 4075424, 4075427,
	3630357, 59542, 4039054, 3641966, 4086032, 3469470, 3469472, 3974461, 3671790, 3659030, 3919342, 3686189,
	3954081, 3954082, 4015598, 4025352, 3659794, 3659795, 3659409, 3659420, 4073499, 4073500, 3697875, 3448732,
	3658364, 3680933, 3421766, 4005427, 4070215, 4053725, 4051174, 3672471, 3605229, 3686187, 3686350, 3624247,
	3644972, 4072281, 4007659, 3691130, 3679968, 3673191, 3694525, 3600023, 3688487, 3688486, 3689731, 3686126,
	3425057, 3425058, 3422551, 2141096, 3445601, 3445602, 3445598, 3694520, 4006039, 14799, 14793, 3501231, 3478735,
	4039025, 3694521, 2112500,
}

func main40()  {
	var tableData [][]string
	var notFoundtableData [][]string
	u := utils.Utils{}


	header := []string{"№","Url", "PLU_ID", "Status"}

	for Istep, productId := range productIds {
		Istep += 1

		fakeTask, _ := fake.CreateFakeTask(interfaces.FakeTaskOptions{
			UserAgent:   false,
			Cookies:     false,
			Auth:        false,
			Proxy:       false,
			AllowedHttp: nil,
		})

		fakeTask.Fetcher.Auth = map[string]string{
			"password": "sUwF}r#LXcly8%U5",
			"login": "Esom",
		}

		fakeTask.Processor.Save.Project = map[string]string{
			"id": fmt.Sprintf("%d",productId),
		}


		fakeTask.URL = "https://httpstat.us/200"
		fakeTask.URL = fmt.Sprintf("http://media.x5.ru/rest/x5/esom?plu=%d", productId)
		//fakeTask.URL = fmt.Sprintf("http://192.168.41.242:8000/api/v1/data/get_photo?id=%d", productId)

		//println(fakeTask.URL)
		//u := utils.Utils{}
		//println(u.PrintPrettyJson(fakeTask))


		//pipeline := httpPipeline.Constructor(fakeTask, nil)
		//
		//err, _ := pipeline.Run(context.TODO())
		//if err != nil {
		//	return
		//}



		//if fakeTask.Fetcher.Response.Code == 200 {
		//
		//	_ = u.DumpToFile(&interfaces.FileObject{
		//		Path:       fmt.Sprintf("images/%d.jpg", productId),
		//		Data:       string(responseData.([]byte)),
		//	})
		//
		//	tableData = append(tableData, []string{
		//		strconv.Itoa(Istep),
		//		fakeTask.URL,
		//		strconv.Itoa(productId),
		//		strconv.Itoa(fakeTask.Fetcher.Response.Code),
		//	})
		//} else {
		//	notFoundtableData = append(notFoundtableData, []string{
		//		strconv.Itoa(Istep),
		//		fakeTask.URL,
		//		strconv.Itoa(productId),
		//		strconv.Itoa(fakeTask.Fetcher.Response.Code),
		//	})
		//}



		//break
	}

	color.Yellow("table length :%d, notFoundtableData: %d", len(tableData), len(notFoundtableData))
	u.RenderTableOutput(header, tableData)

	u.RenderTableOutput(header, notFoundtableData)
}