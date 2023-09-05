package main

import (
	"log"

	Gin "github.com/gin-gonic/gin"
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/extensions/servers/gin"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/resources/home"
	"github.com/vortex14/gotyphoon/interfaces"

	LOG "github.com/vortex14/gotyphoon/log"
)

type MeParams struct {
	Id string `form:"id" binding:"required" json:"id"`
	//Name string `uri:"name" binding:"required"`
}

//type ProductCreate struct {
//	Name  *string `json:"name" binding:"required"`
//	Price *int    `json:"price" binding:"required"`
//}
//
//type Header struct {
//	UserId *int `header:"user-id" binding:"required"`
//}
//
//func main() {
//	r := gin.Default()
//
//	r.POST("/product", func(c *gin.Context) {
//		data := &ProductCreate{}
//		header := &Header{}
//
//		// bind the headers to data
//		if err := c.ShouldBindHeader(header); err != nil {
//			c.JSON(400, struct {
//				Message string
//				Status  bool
//			}{
//				Message: "пиздец",
//				Status:  false,
//			})
//			return
//		}
//
//		// bind the body to data
//		if err := c.ShouldBindJSON(data); err != nil {
//			c.JSON(400, err.Error())
//			return
//		}
//
//		c.JSON(200, data)
//	})
//
//	r.Run(":8083")
//}

func init() {
	LOG.InitP(false, "DEBUG")
}

type ErrorResponseHeaderModel struct {
	Message string
	Data    string
}

type HeaderModel struct {
	UserId int    `header:"user-id" binding:"required" json:"user-id" description:"Это юзер хедер"`
	Test   string `binding:"required" json:"test" description:"Another"`
}

var TestAction = &gin.Action{
	Action: &forms.Action{
		//Params: &MeParams{},
		Headers: forms.HeaderRequestModel{
			ErrorModel: ErrorResponseHeaderModel{
				Message: "header error",
				Data:    "difficult error",
			},
			ErrorStatusCode: 403,
			Model:           &HeaderModel{},
		},
		MetaInfo: &label.MetaInfo{
			Path:        "headers",
			Name:        "HEADERS",
			Description: "HEADERS DESC",
			Tags:        []string{"headers"},
		},
		Methods: []string{interfaces.GET},
	},
	GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
		ctx.JSON(200, struct {
			Message string
			Code    int
		}{
			Message: "test",
			Code:    200,
		})
	},
}

func Constructor() interfaces.ResourceInterface {
	return home.Constructor("/").
		AddAction(TestAction)
}

func main() {
	server := (&gin.TyphoonGinServer{
		TyphoonServer: &forms.TyphoonServer{
			ActiveSwagger: true,
			Host:          "localhost",
			Port:          80,
			Schema:        "http",
			Level:         interfaces.DEBUG,
			MetaInfo: &label.MetaInfo{
				Name:        "test header",
				Description: "test header desc",
				Version:     "1.0.0",
			},
		},
	}).Init()

	server.AddResource(Constructor())

	//server.AddResource(&forms.Resource{
	//	MetaInfo: &label.MetaInfo{
	//		Path:        "headers",
	//		Name:        "headers",
	//		Description: "header description",
	//	},
	//	Actions: map[string]interfaces.ActionInterface{
	//		"check": &gin.Action{
	//			Action: &forms.Action{
	//				Path:    "check",
	//				Methods: []string{"GET"},
	//				MetaInfo: &label.MetaInfo{
	//					Name:        "check",
	//					Description: "sdfsf",
	//				},
	//
	//			},
	//			GinController: func(ctx *Gin.Context, logger interfaces.LoggerInterface) {
	//				ctx.JSON(200, &struct {
	//					Message string
	//				}{
	//					Message: "sfsdfsdfsdf",
	//				})
	//			},
	//		},
	//	},
	//})

	log.Fatal(server.Run())
}
