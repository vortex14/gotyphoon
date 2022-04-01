package main

import (
	"github.com/fatih/color"
	"github.com/vortex14/gotyphoon/elements/models/timer"
	. "github.com/vortex14/gotyphoon/extensions/models/cmd"
)

func main() {
	cmd := &Command{
		Cmd:  "python3.8",
		Dir:  "./main/project/python3/tp",
		Args: []string{"processor.py", "--config=config.local.yaml", "--level=DEBUG"},
	}
	err := cmd.Run()

	if err != nil {
		color.Red(err.Error())
		return
	}

	go func() {
		for it := range cmd.Output {
			println(it)
		}
	}()

	go func() {
		for it := range cmd.OutputErr {
			println(it)
		}
	}()

	timer.SetTimeout(func(args ...interface{}) {
		cmd.Close()
	}, 1000*10)

	color.Yellow(">>>>")
	cmd.Await()
	color.Red("DONE !")
}
