package fakes

import (
	"github.com/vortex14/gotyphoon/elements/forms"
	"github.com/vortex14/gotyphoon/elements/models/label"
	"github.com/vortex14/gotyphoon/extensions/servers/gin"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/controllers/graph"
	"github.com/vortex14/gotyphoon/extensions/servers/gin/controllers/ping"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

const (
	NAME               = "Fakes"
	PATH               = "/"
	WATERMARK          = "image.typhoon.dev"
	DESCRIPTION        = "Server for data fakes"

	ResourceName       = "data fakers resource"

	FakeUPCPath        = "upc"
	FakeTaskPath       = "task"
	FakeProxyPath      = "proxy"
	FakeImagePath      = "image"
	FakeChargePath     = "charge"
	FakeProductPath    = "product"
	FakePaymentPath    = "payment"
	FakeShippingPath   = "shipping"
	FakeCustomerPath   = "customer"
	FakeCategoryPath   = "category"
	FakeCategoriesPath = "categories"

)

func Constructor(
	port int,

	tracingOptions *interfaces.TracingOptions,
	loggerOptions *log.Options,
	swaggerOptions *interfaces.SwaggerOptions,

) interfaces.ServerInterface {
	return (
		&gin.TyphoonGinServer{
			TyphoonServer: &forms.TyphoonServer{
				Port: port,
				Level: interfaces.INFO,
				MetaInfo: &label.MetaInfo{
					Name        : NAME,
					Description : DESCRIPTION,
				},
				TracingOptions  : tracingOptions,
				LoggerOptions   : loggerOptions,
				SwaggerOptions  : swaggerOptions,
			},
		}).
		Init().
		InitLogger().
		AddResource(
			&forms.Resource{
				MetaInfo: &label.MetaInfo{
					Path: PATH,
					Name: ResourceName,
					Description: DESCRIPTION,
				},
				Actions: map[string]interfaces.ActionInterface{
					ping.PATH          : ping.Controller,
					graph.PATH         : graph.Controller,
					FakeUPCPath        : CreateUpcAction(),
					FakeTaskPath       : CreateTaskAction(),
					FakeProxyPath      : CreateProxyAction(),
					FakeImagePath      : CreateImageAction(),
					FakeChargePath     : CreateChargeAction(),
					FakeProductPath    : CrateProductAction(),
					FakePaymentPath    : CreatePaymentAction(),
					FakeCategoryPath   : CreateCategoryAction(),
					FakeCustomerPath   : CreateCustomerAction(),
					FakeShippingPath   : CreateShippingAction(),
					FakeCategoriesPath : CreateCategoriesAction(),
				},
			},
		)
	}