package python3

// /* ignore for building amd64-linux

import Python3 "github.com/DataDog/go-python3"

type PythonList struct {
	List *Python3.PyObject
	Length int
}


func (l *PythonList) GetIter() *Python3.PyObject {
	return l.List.GetIter()
}

func (l *PythonList) GetLength() int {
	return Python3.PyList_Size(l.List)
}

// */
