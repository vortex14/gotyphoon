package code

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)


func TestCommentCode(t *testing.T) {


	Convey("init test code", t, func() {
			testCode := `
package interfaces
//myapp
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

//my comment


///* ignore for building amd64-linux
var ignoreVariable str
// */

type PipelineGroup struct {
	
	///* ignore for building amd64-linux
	test str
	//*/
	
	*label.MetaInfo

	LambdaMap     map[string]interfaces.LambdaInterface  //test comment
	PyLambdaMap   map[string]interfaces.LambdaInterface  //test 2 comment

	Stages        []interfaces.BasePipelineInterface
	Consumers     map[string]interfaces.ConsumerInterface

	graph         interfaces.GraphInterface
	LOG           interfaces.LoggerInterface

}


type Action struct {

	PyController   interfaces.Controller  //Python 
	Middlewares    [] interfaces.MiddlewareInterface  //Before

}

var c int

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

	// */
// comment
}


`
			testcommented := `
package interfaces
//myapp
import (
    "os"
///* ignore for building amd64-linux
//    "fmt"
// 
//*/
    "test/lib"
)

var i int
// /* ignore for building amd64-linux
//func (g *PipelineGroup) SetLogger(logger interfaces.LoggerInterface)  {
//	g.LOG = logger
//}
//*/
var a int

//my comment


///* ignore for building amd64-linux
//var ignoreVariable str
// */

type PipelineGroup struct {
	
	///* ignore for building amd64-linux
//	test str
	//*/
	
	*label.MetaInfo

	LambdaMap     map[string]interfaces.LambdaInterface  //test comment
	PyLambdaMap   map[string]interfaces.LambdaInterface  //test 2 comment

	Stages        []interfaces.BasePipelineInterface
	Consumers     map[string]interfaces.ConsumerInterface

	graph         interfaces.GraphInterface
	LOG           interfaces.LoggerInterface

}


type Action struct {

	PyController   interfaces.Controller  //Python 
	Middlewares    [] interfaces.MiddlewareInterface  //Before

}

var c int

type PipelineGroupGraph interface {


	// /* ignore for building amd64-linux
//
//	SetGraph(graph GraphInterface)
//	InitGraph(parentNode string)
//	SetGraphNodes(nodes map[string]NodeInterface)
//
	// */

	PipelineGroupInterface
}

type PipelineGroupInterface interface {
	Run(ctx context.Context)
	GetName() string
	GetFirstPipelineName() string
	SetLogger(logger LoggerInterface)

	// /* ignore for building amd64-linux
//	SetGraph(graph GraphInterface)
//
//	InitGraph(parentNode string)
//
//	SetGraphNodes(nodes map[string]NodeInterface)
//
	// */
// comment
}


`

		Convey("Run comment code",func() {

			matchCode := "ignore for building amd64-linux"

			data := CommentCode(fmt.Sprintf("/* %s", matchCode), testCode)
			So(testcommented, ShouldContainSubstring, data)

		})

		Convey("Run uncomment code",func() {

			matchCode := "ignore for building amd64-linux"

			data := UnCommentCode(fmt.Sprintf("/* %s", matchCode), testcommented)
			So(data, ShouldContainSubstring, testCode)
		})



	})




}
