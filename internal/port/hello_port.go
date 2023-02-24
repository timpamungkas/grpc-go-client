package port

import (
	"context"
)

type HelloClientPort interface {
	SayHello(ctx context.Context, name string) (string, error)
}
