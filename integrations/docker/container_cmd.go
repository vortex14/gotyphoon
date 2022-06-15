package docker

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/docker/docker/client"
	"github.com/vortex14/gotyphoon/elements/models/awaitabler"
	"github.com/vortex14/gotyphoon/elements/models/singleton"
	Errors "github.com/vortex14/gotyphoon/errors"
)

// Grabbed idea from https://github.com/ahmetb/go-dexec, because used not official docker client.

// Execution determines how the command is going to be executed. Currently,
// the only method is ByCreatingContainer.
type Execution interface {
	create(d *client.Client, cmd []string) error
	run(d *client.Client, stdin io.Reader, stdout, stderr io.Writer) error
	wait(d *client.Client) (int, error)

	setEnv(env []string) error
	setDir(dir string) error
}

type ContainerCMD struct {

	singleton.Singleton
	awaitabler.Object

	// Method provides the execution strategy for the context of the Cmd.
	// An instance of Method should not be reused between Cmds.
	Method Execution

	// Path is the path or name of the command in the container.
	Path string

	// Arguments to the command in the container, excluding the command
	// name as the first argument.
	Args []string

	// Env is environment variables to the command. If Env is nil, Run will use
	// Env specified on Method or pre-built container image.
	Env []string

	// Dir specifies the working directory of the command. If Dir is the empty
	// string, Run uses Dir specified on Method or pre-built container image.
	Dir string

	// Stdin specifies the process's standard input.
	// If Stdin is nil, the process reads from the null device (os.DevNull).
	//
	// Run will not close the underlying handle if the Reader is an *os.File
	// differently than os/exec.
	Stdin io.Reader

	// Stdout and Stderr specify the process's standard output and error.
	// If either is nil, they will be redirected to the null device (os.DevNull).
	//
	// Run will not close the underlying handles if they are *os.File differently
	// than os/exec.
	Stdout io.Writer
	Stderr io.Writer

	DockerClient   *client.Client
	started        bool
	closeAfterWait []io.Closer
}

// NewCommand Command returns the Cmd struct to execute the named program with given
// arguments using specified execution method.
//
// For each new Cmd, you should create a new instance for "method" argument.
func NewCommand(client *client.Client, method Execution, name string, arg ...string) *ContainerCMD {
	return &ContainerCMD{
		Method: method,
		Path: name,
		Args: arg,
		DockerClient: client,
	}
}


// Start starts the specified command but does not wait for it to complete.
func (c *ContainerCMD) Start() error {
	if c.Dir != ""  { if err := c.Method.setDir(c.Dir); err != nil { return err }}

	if c.Env != nil { if err := c.Method.setEnv(c.Env); err != nil { return err }}

	if c.started { return Errors.DockerCommandAlreadyStarted }

	c.started = true

	if c.Stdin == nil  { c.Stdin = empty }
	if c.Stdout == nil { c.Stdout = ioutil.Discard }
	if c.Stderr == nil { c.Stderr = ioutil.Discard }

	cmd := append([]string{c.Path}, c.Args...)
	if err := c.Method.create(c.DockerClient, cmd); err != nil { return err }
	if err := c.Method.run(c.DockerClient, c.Stdin, c.Stdout, c.Stderr); err != nil { return err }
	return nil
}

// Wait waits for the command to exit. It must have been started by Start.
//
// If the container exits with a non-zero exit code, the error is of type
// *ExitError. Other error types may be returned for I/O problems and such.
//
// Different than os/exec.Wait, this method will not release any resources
// associated with Cmd (such as file handles).
func (c *ContainerCMD) Wait() error {
	defer closeFds(c.closeAfterWait)
	if !c.started { return Errors.DockerCommandNotStarted }
	ec, err := c.Method.wait(c.DockerClient)

	if err != nil { return err }
	if ec != 0 {
		return &ExitError{ExitCode: ec}
	}
	return nil
}

// Run starts the specified command and waits for it to complete.
//
// If the command runs successfully and copying streams are done as expected,
// the error is nil.
//
// If the container exits with a non-zero exit code, the error is of type
// *ExitError. Other error types may be returned for I/O problems and such.
func (c *ContainerCMD) Run() error {
	if err := c.Start(); err != nil {
		return err
	}
	return c.Wait()
}

// CombinedOutput runs the command and returns its combined standard output and
// standard error.
//
// Docker API does not have strong guarantees over ordering of messages. For instance:
//     >&1 echo out; >&2 echo err
// may result in "out\nerr\n" as well as "err\nout\n" from this method.
func (c *ContainerCMD) CombinedOutput() ([]byte, error) {
	if c.Stdout != nil { return nil, Errors.DockerStoutAlreadySet  }
	if c.Stderr != nil { return nil, Errors.DockerStdErrAlreadySet }
	var b bytes.Buffer
	c.Stdout, c.Stderr = &b, &b
	err := c.Run()
	return b.Bytes(), err
}

// Output runs the command and returns its standard output.
//
// If the container exits with a non-zero exit code, the error is of type
// *ExitError. Other error types may be returned for I/O problems and such.
//
// If c.Stderr was nil, Output populates ExitError.Stderr.
func (c *ContainerCMD) Output() ([]byte, error) {
	if c.Stdout != nil { return nil, Errors.DockerStoutAlreadySet }
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout

	captureErr := c.Stderr == nil
	if captureErr {
		c.Stderr = &stderr
	}
	err := c.Run()
	if err != nil && captureErr {
		if ee, ok := err.(*ExitError); ok {
			ee.Stderr = stderr.Bytes()
		}
	}
	return stdout.Bytes(), err
}

// StdinPipe returns a pipe that will be connected to the command's standard input
// when the command starts.
//
// Different than os/exec.StdinPipe, returned io.WriteCloser should be closed by user.
func (c *ContainerCMD) StdinPipe() (io.WriteCloser, error) {
	if c.Stdin != nil { return nil, Errors.DockerStdInAlreadySet }
	pr, pw := io.Pipe()
	c.Stdin = pr
	return pw, nil
}

// StdoutPipe returns a pipe that will be connected to the command's standard output when
// the command starts.
//
// Wait will close the pipe after seeing the command exit or in error conditions.
func (c *ContainerCMD) StdoutPipe() (io.ReadCloser, error) {
	if c.Stdout != nil { return nil, Errors.DockerStoutAlreadySet }
	pr, pw := io.Pipe()
	c.Stdout = pw
	c.closeAfterWait = append(c.closeAfterWait, pw)
	return pr, nil
}

// StderrPipe returns a pipe that will be connected to the command's standard error when
// the command starts.
//
// Wait will close the pipe after seeing the command exit or in error conditions.
func (c *ContainerCMD) StderrPipe() (io.ReadCloser, error) {
	if c.Stderr != nil { return nil, Errors.DockerStdErrAlreadySet }
	pr, pw := io.Pipe()
	c.Stderr = pw
	c.closeAfterWait = append(c.closeAfterWait, pw)
	return pr, nil
}
