package hello

import (
	"context"

	pb "github.com/timpamungkas/course-grpc-proto/protogen/go/hello"
	"google.golang.org/grpc"
)

type HelloAdapter struct {
	helloClient pb.HelloServiceClient
}

func NewHelloAdapter(conn *grpc.ClientConn) (*HelloAdapter, error) {
	client := pb.NewHelloServiceClient(conn)

	return &HelloAdapter{
		helloClient: client,
	}, nil
}

func (a *HelloAdapter) SayHello(ctx context.Context, name string) (string, error) {
	helloRequest := &pb.HelloRequest{
		Name: name,
	}

	greet, err := a.helloClient.SayHello(ctx, helloRequest)

	if err != nil {
		return "", err
	}

	return greet.Greet, nil
}
