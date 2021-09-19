package p

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/vortex14/gotyphoon/utils"
	"testing"
)


func TestCommentCode(t *testing.T) {

	Convey("init test code", t, func() {


		testCode := `
package interfaces

import (
    "os"
///* ignore for building amd64-linux
    "fmt"
 
//*/
    "test/lib"
)

var i int
// /* ignore for building amd64-linux
func (g *PipelineGroup) SetLogger(logger interfaces.LoggerInterface)  {
	g.LOG = logger
}
//*/
var a int

// my comment

func TestFunction(x int) {
	println("x: ",x)
}

///* ignore for building amd64-linux
var ignoreVariable str
// */

type PipelineGroup struct {
	
	///* ignore for building amd64-linux
	test str
	//*/
	
	*label.MetaInfo

	LambdaMap     map[string]interfaces.LambdaInterface // test comment
	PyLambdaMap   map[string]interfaces.LambdaInterface // test 2 comment

	Stages        []interfaces.BasePipelineInterface
	Consumers     map[string]interfaces.ConsumerInterface

	graph         interfaces.GraphInterface
	LOG           interfaces.LoggerInterface

}

var c int


//
import (
	"context"
)

type BasePipelineInterface interface {
	Run(
		context context.Context,
		reject func(pipeline BasePipelineInterface, err error),
		next func(ctx context.Context),
	)
	RunMiddlewareStack(
		context context.Context,
		reject func(middleware MiddlewareInterface, err error),
		next func(ctx context.Context),
	)
	Cancel(
		context context.Context,
		logger LoggerInterface,
		err error,
	)
	MetaDataInterface
}

type ProcessorPipelineInterface interface {
	BasePipelineInterface
	Crawl()
	Finish()
	Switch()
}

type PipelineGroupGraph interface {


	// /* ignore for building amd64-linux

	SetGraph(graph GraphInterface)
	InitGraph(parentNode string)
	SetGraphNodes(nodes map[string]NodeInterface)

	// */

	PipelineGroupInterface
}

type PipelineGroupInterface interface {
	Run(ctx context.Context)
	GetName() string
	GetFirstPipelineName() string
	SetLogger(logger LoggerInterface)

	// /* ignore for building amd64-linux
	SetGraph(graph GraphInterface)

	InitGraph(parentNode string)

	SetGraphNodes(nodes map[string]NodeInterface)
//
	// */

}


type CallbackPipelineInterface interface {
	Call(ctx context.Context, data interface{})
}



type ConsumerInterface interface {

}

type LambdaInterface interface {

}

type HandlerInterface interface {
	
}

type ResponseInterface interface {
	
}


`

		matchCode := "/* ignore for building amd64-linux"

		//excludeDirs := map[string]bool{"vendor": true, ".git": true, "tmp": true, ".idea": true}
		//UncommentDir(startDir, matchCode, excludeDirs)
		//utils.CommentDir(startDir, matchCode, excludeDirs)
		data := utils.CommentCode(matchCode, testCode)
		println(data)



	})

}
