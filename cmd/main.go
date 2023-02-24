package main

import (
	"context"
	"log"

	"github.com/timpamungkas/grpc-go-client/internal/adapter/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(logWriter{})

	var opts []grpc.DialOption
	opts = append(opts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	conn, err := grpc.Dial("localhost:9090", opts...)
	if err != nil {
		log.Fatalln("Can't connect to gRPC server : ", err)
	}

	defer conn.Close()

	helloAdapter, err := hello.NewHelloAdapter(conn)
	if err != nil {
		log.Fatalf("Failed to initialize hello adapter : %v\n", err)
	}

	greet, err := helloAdapter.SayHello(context.Background(), "Bruce Wayne")

	if err != nil {
		log.Fatalf("Failed to call SayHello : %v\n", err)
	}

	log.Println(greet)
}
