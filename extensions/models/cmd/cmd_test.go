package cmd

import (
	"github.com/fatih/color"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
	"time"
)

func TestRunLsCMD(t *testing.T) {
	Convey("cmd ls", t, func(c C) {

		cmd := Command{
			Cmd:  "ls",
			Dir:  ".",
			Args: []string{"-la"},
		}

		err := cmd.Run()
		So(err, ShouldBeNil)

		go func() {
			var files []string

			for it := range cmd.Output {

				switch {
				case strings.Contains(it, "cmd.go"):
					c.So(true, ShouldBeTrue)
					files = append([]string{it}, files...)
				case strings.Contains(it, "cmd_test.go"):
					c.So(true, ShouldBeTrue)
					files = append([]string{it}, files...)

				}
			}

			c.So(len(files), ShouldEqual, 2)

		}()

		go func() {
			time.Sleep(2 * time.Second)
			cmd.Close()
		}()
		//
		time.Sleep(7 * time.Second)
		cmd.Await()
		So(true, ShouldBeTrue)

	})
}

func TestRunAwaitCmd(t *testing.T) {
	Convey("sleep", t, func(c C) {
		cmd := Command{Cmd: "sleep", Args: []string{"5"}}
		st := cmd.RunAwait()
		//println(st)
		color.Yellow("%+v", st)
		c.So(st, ShouldNotBeNil)
		c.So(st.Complete, ShouldBeTrue)

	})
}
