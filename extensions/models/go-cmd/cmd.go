package go_cmd

import (
	"sync"
	"time"

	"github.com/go-cmd/cmd"

	"github.com/vortex14/gotyphoon/elements/models/awaitable"
	"github.com/vortex14/gotyphoon/elements/models/singleton"
	Errors "github.com/vortex14/gotyphoon/errors"
)

type Command struct {
	singleton.Singleton
	awaitable.Object
	mu sync.Mutex

	Cmd string
	isDone bool
	Refresh float32 // sec
	cmd *cmd.Cmd
	Args []string
	countRead int
	Output chan string

}

func (c *Command) init()  {
	c.Construct(func() {
		c.Output = make(chan string)
		useCmd := cmd.NewCmd(c.Cmd, c.Args...)
		useCmd.Start()
		c.cmd = useCmd
		c.Add()
		go c.checkingOutput()

	})
}

func (c *Command) checkingOutput()  {
	ticker := time.NewTicker(time.Duration(c.Refresh * 1000) * time.Millisecond)
	for range ticker.C {
		status := c.cmd.Status()
		n := len(status.Stdout)
		countNeed := n - c.countRead

		var startRead int
		if countNeed == n { startRead = 0 }
		if c.countRead > 0  { startRead = c.countRead }

		iterations := countNeed + c.countRead
		c.printCmdOutput(&status, startRead, iterations)
		if status.Complete && !c.isDone{ c.close(); return }

	}
}

func (c *Command) close()  {
	c.mu.Lock()
	c.Done()
	c.isDone = true
	close(c.Output)
	c.mu.Unlock()
}

func (c *Command) printCmdOutput(status *cmd.Status, startRead int, iterations int)  {
	for i := startRead; i < iterations; i++ {
		//color.Yellow(status.Stdout[i])
		c.Output <- status.Stdout[i]
		//color.Yellow(fmt.Sprintf("%+v", status))
		//fmt.Println("read row: ", i , "row: " ,status.Stdout[i], "stdout len: " ,len(status.Stdout))
		c.countRead ++
	}
}

func (c *Command) Run() error {
	if len(c.Cmd) == 0 { return Errors.ErrorCmdNotFound}
	c.init()
	return nil
}