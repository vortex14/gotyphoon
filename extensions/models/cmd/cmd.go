package cmd

import (
	"sync"

	"github.com/fatih/color"
	"github.com/go-cmd/cmd"
	"github.com/hashicorp/go-multierror"

	"github.com/vortex14/gotyphoon/elements/models/awaitabler"
	"github.com/vortex14/gotyphoon/elements/models/singleton"
	Errors "github.com/vortex14/gotyphoon/errors"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/log"
)

// Example from https://github.com/go-cmd/cmd/blob/master/examples/blocking-streaming/main.go

func init()  {
	log.InitD()
}

type Command struct {
	LOG interfaces.LoggerInterface
	singleton.Singleton
	awaitabler.Object
	mu sync.Mutex
	cmd *cmd.Cmd

	Cmd string
	isDone bool
	Dir  string
	Errors error
	Args []string
	status chan bool
	Output chan string
	OutputErr chan string
}

func (c *Command) init()  {
	c.Construct(func() {
		c.LOG = log.New(log.D{"cmd": c.Args[0]})
		c.Output = make(chan string)
		c.OutputErr = make(chan string)

		cmdOptions := cmd.Options{ Buffered: false, Streaming: true }

		useCmd := cmd.NewCmdOptions(cmdOptions, c.Cmd, c.Args...)
		useCmd.Start()
		useCmd.Status()

		if len(c.Dir) > 0 { useCmd.Dir = c.Dir }

		c.cmd = useCmd
		c.Add()
		go c.readOutputStream()


	})
}

func (c *Command) readOutputStream()  {
	c.LOG.Debug("tail -f cmd output")
	for {
		select {
		case line, open := <-c.cmd.Stdout:
			if !open { continue }

			c.Output <- line

		case line, open := <-c.cmd.Stderr:
			if !open { continue }
			c.OutputErr <- line

		case status, ok := <-c.status:
			if !ok || !status {

				err := c.cmd.Stop()

				if err != nil {
					c.Errors = multierror.Append(c.Errors, Errors.ErrorStopCmd)
					color.Red("%s: %s", Errors.ErrorStopCmd.Error(), err.Error())
				}

				return

			}

		}

	}
}

func (c *Command) Close()  {
	c.mu.Lock()
	c.Done()
	c.isDone = true
	c.status <- true
	close(c.Output)
	c.mu.Unlock()
}

func (c *Command) Run() error {
	if len(c.Cmd) == 0 { return Errors.ErrorCmdNotFound}
	c.init()
	return nil
}