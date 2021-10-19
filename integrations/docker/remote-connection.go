package docker

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"io"
	//docker "github.com/fsouza/go-dockerclient"
	//"log"
	"net/http"
	"os"

	"github.com/docker/cli/cli/connhelper"

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
	var clientConnection *client.Client
	var err error
	if d.remoteDockerClient == nil && len(d.RemoteSSHUrl) > 0 {
		err, clientConnection = d.RemoteConnect(); if err != nil { return err, nil }
	} else {
		clientConnection = d.GetClient()
	}
	println("!!! >> > >> >> > > > > >", d.client)

	containers, err := clientConnection.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil { return err, nil }

	for _, container := range containers {
		fmt.Println(container.ID, container.Names[0])
	}

	return nil, containers
}

//https://github.com/moby/moby/blob/8e610b2b55bfd1bfa9436ab110d311f5e8a74dcb/integration/internal/container/exec.go#L38

func (d *Docker) RunSyncRemoteCommandInContainer(context context.Context, container types.Container) (error, ExecResult)  {

	exec, err := Exec(context, d.remoteDockerClient, container.ID, []string{"ls", "/opt/bitnami/solr/data"})
	if err != nil {
		return err, ExecResult{}
	}

	color.Yellow("%+v", exec.outBuffer.String())
	return nil, ExecResult{}
}