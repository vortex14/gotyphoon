package python3

// /* ignore for building amd64-linux
import (
	"log"
	"os"

	Python3 "github.com/DataDog/go-python3"
	"github.com/fatih/color"
)

type PythonModule struct {
	Name      string
	Path      string
	dict      *Python3.PyObject
	module    *Python3.PyObject
	importRef *Python3.PyObject
}

func (m *PythonModule) GetName() string {
	return m.Name
}

func (m *PythonModule) GetPath() string {
	return m.Path
}

func (m *PythonModule) Reset() {

}

func (m *PythonModule) Clean() {
	m.importRef.DecRef()
	m.module.DecRef()
}

func (m *PythonModule) SetImport() {

}

func (m *PythonModule) CreatePythonCallback() {

}

func (m *PythonModule) GetFunction(name string) *Python3.PyObject {
	if m.dict == nil {
		color.Red("Module did not initialized")
		return nil
	}
	function := Python3.PyDict_GetItemString(m.dict, name)

	if !(function != nil && Python3.PyCallable_Check(function)) {
		color.Red("could not find function '%s'", name)
		return nil
	}
	return function
}

func (m *PythonModule) InitModuleHere() bool {
	color.Yellow("Init %s", m.Path)
	var status bool
	if len(m.Path) == 0 {
		color.Red("Path Module not found %s", m.Path)
		status = false
		return status
	} else {
		if _, err := os.Stat(m.Path); os.IsNotExist(err) {
			status = false
			color.Red("Module Path not found %s", m.Path)
			return status
		}
	}

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return status
	}

	color.Yellow("%s", path)
	ret := Python3.PyRun_SimpleString("import sys\nsys.path.append(\"" + path + "\")")
	if ret != 0 {
		status = false
		log.Fatalf("error appending '%s' to python3 sys.path", path)
		return status
	}

	oImport := Python3.PyImport_ImportModule(m.Path)
	if !(oImport != nil && Python3.PyErr_Occurred() == nil) {
		Python3.PyErr_Print()
		color.Red("failed to import module %s", m.Path)
		oImport.DecRef()
		return status
	}

	oModule := Python3.PyImport_AddModule(m.Path)

	if !(oModule != nil && Python3.PyErr_Occurred() == nil) {
		Python3.PyErr_Print()
		color.Red("failed to add module '%s'", m.Path)
		oModule.DecRef()
		return status
	}

	oDict := Python3.PyModule_GetDict(oModule)
	if !(oDict != nil && Python3.PyErr_Occurred() == nil) {
		Python3.PyErr_Print()
		color.Red("could not get dict for %s", m.Path)
		return status
	}

	m.importRef = oImport
	m.module = oModule
	m.dict = oDict

	status = true

	return status

}

// */
