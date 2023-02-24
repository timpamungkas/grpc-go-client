package main

import (
	"context"
	"log"

	"github.com/timpamungkas/grpc-go-client/internal/adapter/hello"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(logWriter{})

	helloAdapter, err := hello.NewHelloAdapter("localhost:9090")
	if err != nil {
		log.Fatalf("Failed to initialize hello adapter : %v\n", err)
	}

	// application := app.NewApplication(helloAdapter)

	greet, err := helloAdapter.SayHello(context.Background(), "Tim")

	if err != nil {
		log.Fatalf("Failed to call SayHello : %v\n", err)
	}

	log.Println(greet)
}
