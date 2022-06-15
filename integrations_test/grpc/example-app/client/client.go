package main

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
}