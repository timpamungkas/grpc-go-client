package application

import (
	"context"

	pb "github.com/timpamungkas/course-grpc-proto/protogen/go/hello"
	"github.com/timpamungkas/grpc-go-client/internal/port"
)

type Application struct {
	helloClient port.HelloClientPort
}

func NewApplication(helloClientPort port.HelloClientPort) *Application {
	return &Application{
		helloClient: helloClientPort,
	}
}

func (a *Application) SayHello(ctx context.Context, request *pb.HelloRequest) (string, error) {
	greet, err := a.helloClient.SayHello(ctx, request)

	if err != nil {
		return "", err
	}

	return greet.Greet, nil
}
