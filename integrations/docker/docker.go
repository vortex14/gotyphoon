package docker

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/vortex14/gotyphoon/elements/models/singleton"

	"github.com/vortex14/gotyphoon/environment"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/utils"
)


//echo -n 'LOGIN:PASSWORD' | base64

type Docker struct {
	singleton.Singleton

	//isLatestTag bool
	LatestTag string
	env       *environment.Settings
	Project   interfaces.Project
	client    *client.Client

	RemoteSSHUrl       string
	remoteDockerClient *client.Client
}

func (d *Docker) GetRemoteClient() *client.Client {
	return d.remoteDockerClient
}

func (d *Docker) GetClient() *client.Client {
	switch {
	case d.remoteDockerClient != nil:
		return d.GetRemoteClient()
	default:
		return d.GetLocalClient()
	}
}

var dockerRegistryUserID = ""

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

type Log struct {
	Stream string `json:"stream"`
	Status string `json:"status"`
}

func (d *Docker) print(rd io.Reader) error {
	var lastLine string

	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		lastLine = scanner.Text()
		color.Yellow("%s", lastLine)
		//var log Log
		//json.Unmarshal(scanner.Bytes(), &log)
		//logLine := fmt.Sprintf("%s, %s", strings.ReplaceAll(log.Stream, "\n", ""), log.Status)
		//color.Yellow(logLine)

	}

	errLine := &ErrorLine{}
	err := json.Unmarshal([]byte(lastLine), errLine)
	if err != nil {
		return err
	}
	if errLine.Error != "" {
		return errors.New(errLine.Error)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (d *Docker) GetLastTagName() string {
	currTime := time.Now()
	tagName := fmt.Sprintf("%s:%d-%d-%d", d.env.DockerImages, currTime.Day(), int(currTime.Month()), currTime.Year())
	return tagName
}

func (d *Docker) GetLocalClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	d.client = cli

	return cli
}

func (d *Docker) CreateNewTag(cli *client.Client)  {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	tagName := d.GetLastTagName()
	errTag := cli.ImageTag(ctx, "typhoon-lite:latest", tagName)
	if errTag != nil {
		color.Red("%s", errTag.Error())
	} else {
		color.Green("created new image tag: %s", tagName)
	}
}

func (d *Docker) build(workDir string, options *archive.TarOptions, opts types.ImageBuildOptions) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	color.Green("PatH: %s", workDir)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	tar, err := archive.TarWithOptions(fmt.Sprintf("%s/", workDir), options)
	if err != nil {
		return err
	}

	if err != nil {
		color.Red("%s", err.Error())
		return err
	}

	res, err := cli.ImageBuild(ctx, tar, opts)
	if err != nil {
		color.Red("%s", err.Error())
		return err
	}

	defer res.Body.Close()
	exceptions := d.print(res.Body)

	if exceptions != nil {
		color.Red("%s", exceptions.Error())
	}
	d.CreateNewTag(cli)

	return nil
}

func (d *Docker) initEnv() {
	envSetting := environment.Environment{}
	_, env := envSetting.GetSettings()
	d.env = env
}

func (d *Docker) GetAuthConfig() types.AuthConfig {
	var authConfig = types.AuthConfig{
		Username: d.env.DockerLogin,
		Password: d.env.DockerPassword,
	}

	return authConfig
}

func (d *Docker) PushImage()  {
	d.initEnv()
	color.Yellow("Typhoon docker push to %s...", d.env.DockerHub)
	b, err := strconv.ParseBool(d.LatestTag)
	if err != nil{

		color.Red(err.Error())
		os.Exit(1)
	}
	var lastTagName string

	if b {
		lastTagName = d.GetLastTagName()
	} else {
		lastTagName = "typhoon-lite:latest"
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	var authConfig = types.AuthConfig{
		Username:      d.env.DockerLogin,
		Password:      d.env.DockerPassword,
		ServerAddress: d.env.DockerHub,
	}
	authConfigBytes, _ := json.Marshal(authConfig)
	authConfigEncoded := base64.URLEncoding.EncodeToString(authConfigBytes)

	opts := types.ImagePushOptions{RegistryAuth: authConfigEncoded}
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	rd, err := dockerClient.ImagePush(ctx, dockerRegistryUserID+ lastTagName, opts)
	if err != nil {
		color.Red("%s", err.Error())
		os.Exit(1)
	}

	defer rd.Close()

	exceptions := d.print(rd)

	if exceptions != nil {
		color.Red("%s", exceptions.Error())
	}
	color.Green("%s", lastTagName)
}

func (d *Docker) BuildImage()  {
	color.Yellow("Typhoon docker build ...")
	d.initEnv()

	options := &archive.TarOptions{
		ExcludePatterns: []string{
			"extensions/tests/*",
			".git",
			"chrome",
	}}

	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{dockerRegistryUserID + "typhoon-lite"},
		Remove:     true,
		AuthConfigs: map[string]types.AuthConfig{
			d.env.DockerHub: d.GetAuthConfig(),
		},

	}

	err := d.build(d.env.Path, options, opts)
	if err != nil {
		color.Red(err.Error())
		return
	}


}

func (d *Docker) ListContainers()  {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
}

func (d *Docker) RunComponent(component string) error {
	d.Project.LoadConfig()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	if err != nil {
		//fmt.Println(err.Error())
		//color.Red("Error: %s", err)
		return err
	}

	containerConfig := &container.Config{
		Image: fmt.Sprintf("typhoon-lite-%s", d.Project.GetName()),
		Cmd: []string{"python", "donor.py --config=config.local.yaml --level=DEBUG"},
	}


	_, err = cli.ContainerCreate(ctx, containerConfig, nil, nil, nil, "typhoon")

	if err != nil {
		color.Red("ContainerCreate: %s", err)
		return err
	}

	return nil
}

func (d *Docker) ProjectBuild()  {
	u := utils.Utils{}

	_, dockerFile := u.GetGoTemplate(&interfaces.FileObject{
		Path: "../builders/v1.1",
		Name: "ProjectDockerfile",

	})
	goTemplateDocker := interfaces.GoTemplate{
		Source: dockerFile,
		ExportPath: "Dockerfile",
		Data: map[string]string{
			"TYPHOON_IMAGE": d.Project.GetDockerImageName(),
		},
	}

	err := u.GoRunTemplate(&goTemplateDocker)
	if !err {
		color.Red("creation Dockerfile was fail")
		os.Exit(1)
	}

	color.Green("Dockerfile created!")

	color.Yellow("Typhoon project docker build ...")

	options := &archive.TarOptions{}
	projectConfig := d.Project.LoadConfig()
	opts := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{"typhoon-lite-" + projectConfig.ProjectName},
		Remove:     true,

	}

	_ = d.build(d.Project.GetProjectPath(), options, opts)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}





}

func (d *Docker) RemoveResources()  {
	u := utils.Utils{}
	u.RemoveFiles([]string{
		"Dockerfile",
	})

	color.Green("Removed")
}
