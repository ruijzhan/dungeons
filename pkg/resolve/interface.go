package resolve

import (
	"context"
	"net"
)

type Resolver interface {
	Resolve(context.Context, string) ([]net.IP, error)
}

func New() Resolver {
	return &DefaultResolver{
		l: net.DefaultResolver,
	}
}
