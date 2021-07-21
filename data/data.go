package data

import (
	"github.com/vortex14/gotyphoon/interfaces"
)

type StructData struct {
	Fields []string
}

func (s *StructData) GetFields()  {

}

func TestFunc () interfaces.TestData {
	test := &StructData{
		Fields: []string{"1", "2"},
	}

	return test

}
