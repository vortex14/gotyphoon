package docker

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/fatih/color"
	"io"
	//docker "github.com/fsouza/go-dockerclient"
	//"log"
	"net/http"
	"os"

	"github.com/docker/cli/cli/connhelper"
	MobyContainer "github.com/moby/moby/integration/internal/container"

	//"github.com/ahmetalpbalkan/dexec"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)


func closeFds(l []io.Closer) {
	for _, fd := range l {
		err := fd.Close()
		if err != nil {
			color.Red("closeFds has: ",err.Error())
			return
		}
	}
}

type emptyReader struct{}

func (r *emptyReader) Read(b []byte) (int, error) { return 0, io.EOF }

var empty = &emptyReader{}


// Provided from
//https://gist.github.com/agbaraka/654a218f8ea13b3da8a47d47595f5d05

func (d *Docker) RemoteConnect() (error, *client.Client)  {
	helper, err := connhelper.GetConnectionHelper(d.RemoteSSHUrl)

	if err != nil{ return err, nil }

	httpClient := &http.Client{
		// No tls
		// No proxy
		Transport: &http.Transport{
			DialContext: helper.Dialer,
		},
	}

	var clientOpts []client.Opt

	clientOpts = append(clientOpts,
		client.WithHTTPClient(httpClient),
		client.WithHost(helper.Host),
		client.WithDialContext(helper.Dialer),

	)

	version := os.Getenv("DOCKER_API_VERSION")

	if version != "" {
		clientOpts = append(clientOpts, client.WithVersion(version))
	} else {
		clientOpts = append(clientOpts, client.WithAPIVersionNegotiation())
	}


	cl, err := client.NewClientWithOpts(clientOpts...)

	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}

	d.remoteDockerClient = cl

	return err, cl
}


func (d *Docker) GetRemoteDockerImagesList() (error, []types.ImageSummary) {
	if d.remoteDockerClient == nil { err, _ := d.RemoteConnect(); if err != nil { return err, nil }}

	list, err := d.remoteDockerClient.ImageList(context.Background(), types.ImageListOptions{})
	return 	err, list
}

func (d *Docker) GetRemoteActiveContainersList() (error, []types.Container) {
	if d.remoteDockerClient == nil { err, _ := d.RemoteConnect(); if err != nil { return err, nil }}

	containers, err := d.remoteDockerClient.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil { return err, nil }

	for _, container := range containers {
		fmt.Println(container.ID)
	}

	return nil, containers
}

//https://github.com/moby/moby/blob/8e610b2b55bfd1bfa9436ab110d311f5e8a74dcb/integration/internal/container/exec.go#L38

func (d *Docker) RunRemoteCommandInContainer(context context.Context, container types.Container) (error, MobyContainer.ExecResult)  {
	attach, err := d.remoteDockerClient.ContainerExecAttach(context, container.ID, types.ExecStartCheck{})
	if err != nil {
		color.Red("%+v", err.Error())
		return err, MobyContainer.ExecResult{}
	}

	config :=  types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		Cmd: []string{"ls", "-la"},
	}

	cresp, err := d.remoteDockerClient.ContainerExecCreate(context, container.ID, config)
	if err != nil {
		return err, MobyContainer.ExecResult{}
	}

	execID := cresp.ID

	// run it, with stdout/stderr attached
	aresp, err := d.remoteDockerClient.ContainerExecAttach(context, execID, types.ExecStartCheck{})
	if err != nil {
		return err, MobyContainer.ExecResult{}
	}

	defer aresp.Close()

	// read the output
	var outBuf, errBuf bytes.Buffer
	outputDone := make(chan error)

	go func() {
		// StdCopy demultiplexes the stream into two buffers
		_, err = stdcopy.StdCopy(&outBuf, &errBuf, aresp.Reader)
		outputDone <- err
	}()

	select {
	case err := <-outputDone:
		if err != nil {
			return err, MobyContainer.ExecResult{}
		}
		break

	case <-context.Done():
		return context.Err(), MobyContainer.ExecResult{}
	}

	// get the exit code
	_, err = d.remoteDockerClient.ContainerExecInspect(context, execID)
	if err != nil {
		return err, MobyContainer.ExecResult{}
	}

	return nil, MobyContainer.ExecResult{}




	//ContainerCMD{
	//	Method:         nil,
	//	Path:           "",
	//	Args:           nil,
	//	Dir:            "",
	//	DockerClient:   nil,
	//	started:        false,
	//	closeAfterWait: nil,
	//}
	//
	//cmd := exec.Command("ls", "-la")
	//cmd.Stdin = attach.Reader
	//
	//var out bytes.Buffer
	//cmd.Stdout = &out
	//errm := cmd.Run()
	//if err != nil {
	//	color.Red(errm.Error())
	//}
	//fmt.Printf("in all caps: %q\n", out.String())







	//attach.Conn.

	//attach.Conn.

	color.Red("%+v", attach)
	//client := d.remoteDockerClient.(*client.Client)
	//err := client.Ping()
	//if err != nil {
	//	return
	//}
	//dc := dexec.Docker{}
	//
	////d.remoteDockerClient..
	//m, _ := dexec.ByCreatingContainer(docker.CreateContainerOptions{
	//	Config: &docker.Config{Image: container.Image}})
	//
	//md := d.Command(m, "echo", `I am running inside a container!`)
	//b, err := cmd.Output()
	//if err != nil { log.Fatal(err) }
	//log.Printf("%s", b)
}