package fakes

// /* ignore for building amd64-linux
//
//import (
//	"github.com/vortex14/gotyphoon/elements/forms"
//	"github.com/vortex14/gotyphoon/elements/models/label"
//	"github.com/vortex14/gotyphoon/extensions/servers/gin/controllers/graph"
//	"github.com/vortex14/gotyphoon/extensions/servers/gin/controllers/ping"
//	graph2 "github.com/vortex14/gotyphoon/extensions/servers/gin/graph"
//	"github.com/vortex14/gotyphoon/interfaces"
//	"github.com/vortex14/gotyphoon/log"
//
//	graphFormExt "github.com/vortex14/gotyphoon/extensions/forms/graph"
//)
//
//func GraphConstructor(
//	port int,
//
//	tracingOptions *interfaces.TracingOptions,
//	loggerOptions *log.Options,
//	swaggerOptions *interfaces.SwaggerOptions,
//
//) interfaces.ServerInterface {
//	return (
//		&graph2.TyphoonGraphGinServer{
//			TyphoonServer: &graphFormExt.TyphoonServer{
//				TyphoonServer: &forms.TyphoonServer{
//					Port: port,
//					Level: interfaces.DEBUG,
//					MetaInfo: &label.MetaInfo{
//						Name        : NAME,
//						Description : DESCRIPTION,
//					},
//					TracingOptions  : tracingOptions,
//					LoggerOptions   : loggerOptions,
//					SwaggerOptions  : swaggerOptions,
//				},
//				BuildGraph:    true,
//			},
//		}).
//		Init().
//		InitLogger().
//		AddResource(
//			&graphFormExt.Resource{
//				Resource: &forms.Resource{
//					MetaInfo: &label.MetaInfo{
//						Path: PATH,
//						Name: ResourceName,
//						Description: DESCRIPTION,
//					},
//					Actions: map[string]interfaces.ActionInterface{
//						ping.PATH          : ping.GraphController,
//						graph.PATH         : graph.Controller,
//						FakeUPCPath        : GraphController,
//						//FakeTaskPath       : CreateTaskAction(),
//						//FakeProxyPath      : CreateProxyAction(),
//						FakeImagePath      : CreateImageAction(),
//						//FakeChargePath     : CreateChargeAction(),
//						//FakeProductPath    : CrateProductAction(),
//						//FakePaymentPath    : CreatePaymentAction(),
//						//FakeCategoryPath   : CreateCategoryAction(),
//						//FakeCustomerPath   : CreateCustomerAction(),
//						//FakeShippingPath   : CreateShippingAction(),
//						//FakeCategoriesPath : CreateCategoriesAction(),
//					},
//				},
//			},
//		)
//}
//
//
// */