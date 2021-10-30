package python3

// /* ignore for building amd64-linux
import (
	Python3 "github.com/DataDog/go-python3"
	"github.com/fatih/color"
)

type PythonFunction struct {
	Function *Python3.PyObject
}

// */

//https:pyo3.rs/v0.14.1/ecosystem/async-await.html

// /* ignore for building amd64-linux
func (f *PythonFunction) await(object *Python3.PyObject) *Python3.PyObject {
	aNext := object.GetAttrString("__await__")
	defer aNext.DecRef()
	coro := aNext.CallObject(nil)
	iter := coro.GetAttrString("__next__")

	defer iter.DecRef()

	item := iter.CallObject(nil)
	itemType := item.Type()
	if itemType == nil && !Python3.PyErr_ExceptionMatches(Python3.PyExc_StopIteration) {
		Python3.PyErr_Print()
		color.Red("error getting item type. fn  __next__ for coroutine is not provided")
		return nil
	}

	defer itemType.DecRef()
	defer item.DecRef()

	_, coroResult, _ := Python3.PyErr_Fetch()

	return coroResult
}

func (f *PythonFunction) IsAwaitable(function *Python3.PyObject) bool {
	var status bool
	await := function.GetAttrString("__await__")
	if await != nil && Python3.PyCallable_Check(await) {
		status = true
	}
	defer await.DecRef()
	return status
}


func (f *PythonFunction) CallObject(args *Python3.PyObject) *Python3.PyObject  {
	result := f.Function.CallObject(nil)
	if f.IsAwaitable(result) {
		result = f.await(result)
	}
	return result

}

// */
