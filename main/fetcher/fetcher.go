package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/vortex14/gotyphoon/elements/models/task"
	"github.com/vortex14/gotyphoon/log"
	"github.com/vortex14/gotyphoon/utils"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Meta struct {
	Kind string `yaml:"kind"`
	Component string `yaml:"component"`
	Version string `yaml:"version"`
}

type TableParams struct {
	Headers []string `yaml:"headers"`
}

type TableOutput struct {
	Type string `yaml:"type"`
	Params 	TableParams `yaml:"params"`

}

type TaskTemplate struct {
	Typhoon Meta `yaml:"typhoon"`
	Task *task.TyphoonTask `yaml:"task"`
	Output TableOutput `yaml:"output"`
}

func init()  {
	log.InitD()
}


func main()  {


	path, err := os.Getwd()
	if err != nil {
		return
	}
	path = filepath.Join(path, "main", "fetcher", "task1.yaml")

	yamlFile, _ := ioutil.ReadFile(path)

	var fetcherTemplate TaskTemplate

	err = utils.YamlLoad(&fetcherTemplate, yamlFile)
	if err != nil {
		return
	}
	color.Red("%+v", fetcherTemplate.Task.GetFetcherUrl())
	divideByZero()
	fmt.Println("we survived dividing by zero!")

}

func divideByZero() {
	defer func() {
		if err := recover(); err != nil {
			println("panic occurred:", err)
		}
	}()
	fmt.Println(divide(1, 0))
}

func divide(a, b int) int {
	return a / b
}