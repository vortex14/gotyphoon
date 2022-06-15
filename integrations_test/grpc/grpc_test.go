package main

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	. "github.com/vortex14/gotyphoon/extensions/models/cmd"
	"github.com/vortex14/gotyphoon/interfaces"
	"github.com/vortex14/gotyphoon/utils"
	"os"
	"path/filepath"
	"strings"
	"time"

	"testing"
)

func TestRunGRPCGenerator(t *testing.T) {
	path, err := os.Getwd()

	Convey("generate grpc server", t, func(c C) {

		_ = os.Remove(filepath.Join(path, "example-app", "proto-app", "proto.pb.go"))
		_ = os.Remove(filepath.Join(path, "example-app", "proto-app", "proto_grpc.pb.go"))
		_ = os.Remove(filepath.Join(path, "example-app", "main.go"))
		_ = os.Remove(filepath.Join(path, "example-app", "client", "client"))
		_ = os.Remove(filepath.Join(path, "example-app", "client", "client.go"))
		_ = os.Remove(filepath.Join(path, "example-app", "example-app"))

		cmd := &Command{
			//Cmd: "ls",
			Cmd: "protoc",
			Dir: ".",
			Args: []string{
				"--go_out=.",
				"--go_opt=paths=source_relative",
				"--go-grpc_out=.",
				"--go-grpc_opt=paths=source_relative",
				"example-app/proto-app/proto.proto",
			},
		}

		status := cmd.RunAwait()
		c.So(status.Error, ShouldBeNil)

		codeServer := `package main

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	pb "github.com/vortex14/gotyphoon/integrations_test/grpc/example-app/proto-app"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	test := "{\"status\":true}"

	return &pb.HelloReply{Message: &test}, nil
}

func (s *server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {

	test := "{\"status\":true}"

	return &pb.HelloReply{Message: &test}, nil
}

func (s *server) Dialog(stream pb.Greeter_DialogServer) error {
	color.Yellow("Run stream dialog")
	status := make(chan bool, 1)
	go func() {
		i := 100
		for {
			if i == 0 {
				status <- true
				return
			}
			if err := stream.Send(&pb.Task{Message: fmt.Sprintf("Server: new message body{ %d } ", i)}); err != nil {
				color.Red("%+v", err)
			}
			i -= 1

		}
	}()
	<-status

	color.Green("stream close")

	return nil

}

func main() {
	listener, err := net.Listen("tcp", ":9999")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}`
		u := utils.Utils{}
		err = u.DumpToFile(&interfaces.FileObject{
			Path: filepath.Join(path, "example-app", "main.go"),
			Data: codeServer,
		})
		c.So(err, ShouldBeNil)

	})

	Convey("Build grpc app", t, func(c C) {

		cmd := &Command{Cmd: "go", Args: []string{"build", "."}, Dir: "example-app"}
		status := cmd.RunAwait()
		c.So(status.Error, ShouldBeNil)
		_, err = os.Stat(filepath.Join(path, "example-app", "example-app"))
		c.So(err, ShouldBeNil)
	})

	serverCmd := &Command{Cmd: "./example-app", Args: []string{""}, Dir: "example-app"}
	Convey("run grpc server", t, func(c C) {
		err = serverCmd.Run()
		c.So(err, ShouldBeNil)
		//time.Sleep(2 * time.Second)
		//serverCmd.Close()
	})

	Convey("generate grpc goClient", t, func(c C) {
		goClient := `package main

import (
	"context"
	"fmt"
	"github.com/vortex14/gotyphoon/log"
	"io"

	pb "github.com/vortex14/gotyphoon/integrations_test/grpc/example-app/proto-app"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	log.InitD()
	LOG := log.New(log.D{"client": true})
	conn, err := grpc.Dial("localhost:9999", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		LOG.Debug(fmt.Sprintf("failed to connect: %v", err))
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	answer, err := client.SayHello(context.Background(), &pb.HelloRequest{})
	if err != nil {
		LOG.Debug(fmt.Sprintf("failed to get hello: %v", err))
	}
	LOG.Debug(fmt.Sprintf("say: %v", answer))

	stream, err := client.Dialog(context.Background())

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			// read done.
			return
		}
		if err != nil {
			LOG.Error(fmt.Sprintf("client.RouteChat failed: %v", err))
			break
		}

		LOG.Debug(fmt.Sprintf("Client: Got message  %s", in.Message))
	}
}`
		u := utils.Utils{}
		err = u.DumpToFile(&interfaces.FileObject{
			Path: filepath.Join(path, "example-app", "client", "client.go"),
			Data: goClient,
		})
		c.So(err, ShouldBeNil)

	})

	Convey("Build grpc client", t, func(c C) {

		cmd := &Command{Cmd: "go", Args: []string{"build", "."}, Dir: "example-app/client"}
		status := cmd.RunAwait()
		c.So(status.Error, ShouldBeNil)
		_, err = os.Stat(filepath.Join(path, "example-app", "client", "client"))
		c.So(err, ShouldBeNil)
	})

	goClientCmd := &Command{Cmd: "./client", Args: []string{""}, Dir: "example-app/client"}
	Convey("run grpc goClient", t, func(c C) {
		err := goClientCmd.Run()
		c.So(err, ShouldBeNil)

	})

	Convey("check answer from client", t, func(c C) {

		go func() {
			//log := "Client: Got message  Server: new message body{ 8 } "
			i := 0
			count := 100
			log := ""
			logMessages := 1
			for it := range goClientCmd.Output {
				i += 1

				if i >= 4 {

					log = fmt.Sprintf("Client: Got message  Server: new message body{ %d } ", count)

					if strings.Contains(it, log) {
						//println(log, it)
						//c.So(true, ShouldBeTrue)
						logMessages += 1
					}

					count -= 1
				}

				if logMessages == 100 {
					c.So(true, ShouldBeTrue)
				}

			}
		}()

		go func() {
			for it := range goClientCmd.OutputErr {
				println("??????>>>>", it)
				//c.So(it, ShouldBeNil)
			}
		}()
		time.Sleep(10 * time.Second)

	})
}
