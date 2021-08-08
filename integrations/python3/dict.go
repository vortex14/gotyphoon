package python3

import (
	Python3 "github.com/DataDog/go-python3"
	"github.com/fatih/color"
)

type PythonDict struct {
	Dict *Python3.PyObject
}

func (d *PythonDict) Repr() *Python3.PyObject  {
	return d.Dict.Dir().Repr()
}



func (d *PythonDict) SetDictItems()  {

}

func (d *PythonDict) GetDict()  {

}

func (d *PythonDict) CreateDict()  {

}

func (d *PythonDict) GetCountDictKeys() int  {
	return Python3.PyDict_Size(d.Dict)
}




func (d *PythonDict) GetValue(key string) *Python3.PyObject {
	pyKey := Python3.PyUnicode_FromString(key)
	return d.Dict.GetItem(pyKey)
}

func (d *PythonDict) GetDictKeys() []string {
	var dictKeys []string
	list := &PythonList{List: Python3.PyDict_Keys(d.Dict)}
	length := list.GetLength()
	color.Red("Count keys: %d", length)

	seq := list.GetIter() //ret val: New reference
	if !(seq != nil && Python3.PyErr_Occurred() == nil) {
		Python3.PyErr_Print()
		color.Red("error creating iterator for list")
		return nil
	}
	defer seq.DecRef()

	tNext := seq.GetAttrString("__next__") //ret val: new ref
	if !(tNext != nil && Python3.PyCallable_Check(tNext)) {
		color.Red("iterator has no __next__ function")
		return nil
	}


	defer tNext.DecRef()

	for i := 1; i <= length; i++ {
		item := tNext.CallObject(nil) //ret val: new ref
		if item == nil && Python3.PyErr_Occurred() != nil {
			Python3.PyErr_Print()
			color.Red("error getting next item in sequence")
			return nil
		}
		itemType := item.Type()
		if itemType == nil && Python3.PyErr_Occurred() != nil {
			Python3.PyErr_Print()
			color.Red("error getting item type")
			return nil
		}

		defer itemType.DecRef()
		defer item.DecRef()

		keyName := Python3.PyUnicode_AsUTF8(item)

		dictKeys = append(dictKeys, keyName)
	}



	return dictKeys

}


