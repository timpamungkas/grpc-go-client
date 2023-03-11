package main

import (
	"context"
	"log"

	"github.com/timpamungkas/grpc-go-client/internal/adapter/bank"
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

	// helloAdapter, err := hello.NewHelloAdapter(conn)
	// if err != nil {
	// 	log.Fatalf("Failed to initialize hello adapter : %v\n", err)
	// }

	// runSayHello(helloAdapter, "Bruce Wayne")
	// runSayManyHellos(helloAdapter, "Clark Kent")
	// runSayHelloToEveryone(helloAdapter, []string{"Andy", "Bill", "Christian", "Donny", "Edgar"})
	// runSayHelloContinuous(helloAdapter, []string{"Anna", "Bella", "Carol", "Diana", "Emma"})

	bankAdapter, err := bank.NewBankAdapter(conn)
	if err != nil {
		log.Fatalf("Failed to initialize bank adapter : %v\n", err)
	}

	runGetCurrentBalance(bankAdapter, "7835697001")
	runFetchExchangeRates(bankAdapter, "USD", "IDR")
}

func runSayHello(adapter *hello.HelloAdapter, name string) {
	greet, err := adapter.SayHello(context.Background(), name)

	if err != nil {
		log.Fatalf("Failed to call SayHello : %v\n", err)
	}

	log.Println(greet.Greet)
}

func runSayManyHellos(adapter *hello.HelloAdapter, name string) {
	adapter.SayManyHellos(context.Background(), name)
}

func runSayHelloToEveryone(adapter *hello.HelloAdapter, names []string) {
	adapter.SayHelloToEveryone(context.Background(), names)
}

func runSayHelloContinuous(adapter *hello.HelloAdapter, names []string) {
	adapter.SayHelloContinuous(context.Background(), names)
}

func runGetCurrentBalance(adapter *bank.BankAdapter, acct string) {
	bal, err := adapter.GetCurrentBalance(context.Background(), acct)

	if err != nil {
		log.Fatalf("Failed to call GetCurrentBalance : %v\n", err)
	}

	log.Println(bal)
}

func runFetchExchangeRates(adapter *bank.BankAdapter, fromCur string, toCur string) {
	adapter.FetchExchangeRates(context.Background(), fromCur, toCur)
}
