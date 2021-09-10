package main

import (
	"context"
	"github.com/fogleman/gg"
	"github.com/vortex14/gotyphoon/extensions/data/fake"
	"github.com/vortex14/gotyphoon/interfaces"
	"reflect"
)

type ReflectFunc func(
	context context.Context,
	task interfaces.TaskInterface,
	logger interfaces.LoggerInterface,
	imgCtx *gg.Context,
)

func test(args ...interface{})  {
	for i, arg := range args {
		println(i,reflect.TypeOf(arg).)
		println(i,reflect.TypeOf(arg).)
	}

}


type User struct {
	Name string
}

func main()  {
	//test(&User{})
	test(fake.CreateDefaultTask(), &User{})
}
