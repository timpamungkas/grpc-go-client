package hello

import (
	"context"

	pb "github.com/timpamungkas/course-grpc-proto/protogen/go/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type HelloAdapter struct {
	helloClient pb.HelloServiceClient
}

func NewHelloAdapter(helloServiceUrl string) (*HelloAdapter, error) {
	var opts []grpc.DialOption
	opts = append(opts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	conn, err := grpc.Dial(helloServiceUrl, opts...)
	if err != nil {
		return nil, err
	}

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
