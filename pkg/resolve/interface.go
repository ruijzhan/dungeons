package resolve

import (
	"context"
	"net"

	"github.com/ruijzhan/dungeons/pkg/cache"
)

type Resolver interface {
	Resolve(context.Context, string) ([]net.IP, error)
}

func New() Resolver {
	return &DefaultResolver{
		lookup: net.DefaultResolver,
		cache:  cache.NewTTL(),
	}
}
