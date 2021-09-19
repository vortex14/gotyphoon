package interfaces



// /* ignore for building amd64-linux

import Python3 "github.com/DataDog/go-python3"

// */

const (
	PYTHON 						= "Python"

	TYPHOON2PYTHON2FETCHER 		= "fetcher"
	TYPHOON2PYTHON2PROCESSOR 	= "processor"
	TYPHOON2PYTHON2SCHEDULER 	= "scheduler"
	TYPHOON2PYTHON2TRANSPORTER 	= "result_transporter"
	TYPHOON2PYTHON2DONOR 		= "donor"
)

// /* ignore for building amd64-linux

type PythonFunctionInterface interface {
	CallObject(args *Python3.PyObject) *Python3.PyObject
	IsAwaitable(object *Python3.PyObject) bool
}

type PythonDictInterface interface {
	GetDict()
	CreateDict()
	SetDictItems()
	GetCountDictKeys() int
	GetDictKeys() []string
	GetValue(key string) *Python3.PyObject
}

type PythonListInterface interface {
	GetLength() int
}

type PythonBaseTypeInterface interface {
	CreateString()
	CreateNumber()
}


type PythonClassInterface interface {

}

type PythonTypesInterface interface {
	PythonDictInterface
	PythonBaseTypeInterface
	PythonFunctionInterface
	PythonClassInterface

}

type PythonModuleInterface interface {
	InitModuleHere() bool
	Clean()
	SetImport()
	Reset()
	GetName() string
	GetPath() string

}

type PythonInterface interface {
	InitEnvironment()
	CloseEnvironment()
	InitModules()
}

// */