package python3

///* ignore for building amd64-linux
import (
	"fmt"
	Python3 "github.com/DataDog/go-python3"
	"github.com/fatih/color"
	"github.com/vortex14/gotyphoon/interfaces"
	"os"
)

type Python struct {
	Modules map[string] *PythonModule
}

func (p *Python) InitEnvironment()  {
	Python3.Py_Initialize()
	if !Python3.Py_IsInitialized() {
		fmt.Println("Error initializing the python interpreter")
		os.Exit(1)
	}
}

func (p *Python) RunFile(path string)  bool {
	s, error := Python3.PyRun_AnyFile(path)
	if error != nil {
		color.Red(error.Error())
		panic(error)
	}
	color.Yellow("Running status %d", s)
	return true
}



func (p *Python) InitModules()  {
	for name := range p.Modules {
		module := p.Modules[name]
		status := module.InitModuleHere()
		if status {
			color.Green("Python module %s initialized", name)
		} else {
			color.Red("Python module %s initialization error", name)
		}
	}
}

func (p *Python) Init()  {
	p.InitEnvironment()
	p.InitModules()
}

func (p *Python) CloseEnvironment()  {
	for name := range p.Modules {
		module := p.Modules[name]
		module.Clean()
	}
	Python3.Py_Finalize()
}



func (p *Python) GetInt(object *Python3.PyObject) int {
	return Python3.PyLong_AsLong(object)
}

func ToString(object *Python3.PyObject) string  {
	return Python3.PyUnicode_AsUTF8(object)
}

func (p *Python) GetNewList(length int) interfaces.PythonListInterface {
	newList := Python3.PyList_New(length)
	return &PythonList{List: newList, Length: length}

}


func PyMethodCheck(object *Python3.PyObject) int {
	return 1
}
//*/