package port

import (
	"context"

	pb "github.com/timpamungkas/course-grpc-proto/protogen/go/hello"
	"google.golang.org/grpc"
)

type HelloClientPort interface {
	SayHello(ctx context.Context, request *pb.HelloRequest, opts ...grpc.CallOption) (*pb.HelloResponse, error)
}
