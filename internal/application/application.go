package application

import (
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
