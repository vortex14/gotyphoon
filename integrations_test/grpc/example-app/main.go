package main

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
}