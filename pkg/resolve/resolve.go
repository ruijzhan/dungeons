package resolve

import (
	"context"
	"fmt"
	"net"

	"github.com/ruijzhan/dungeons/pkg/cache"
	"golang.org/x/sync/singleflight"
)

type lookUp interface {
	LookupIPAddr(context.Context, string) ([]net.IPAddr, error)
}

type DefaultResolver struct {
	group  singleflight.Group
	lookup lookUp
	cache  cache.Cache
}

var _ Resolver = (*DefaultResolver)(nil)

func (r *DefaultResolver) Resolve(ctx context.Context, host string) ([]net.IP, error) {
	if r.cache == nil {
		return nil, fmt.Errorf("cache is nil")
	}

	cached, err := r.cache.Get(host)
	if err == nil {
		return cached.([]net.IP), nil
	}

	resolve := func() (interface{}, error) {
		ips, err := r.lookup.LookupIPAddr(ctx, host)
		if err != nil {
			return nil, err
		}
		if len(ips) == 0 {
			return nil, fmt.Errorf("no IP addresses found for host %s", host)
		}

		results := make([]net.IP, len(ips))
		for i, ip := range ips {
			results[i] = ip.IP
		}
		r.cache.Set(host, results)

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			return results, nil
		}

	}

	do, err, _ := r.group.Do(host, resolve)
	if err != nil {
		return nil, err
	}

	return do.([]net.IP), nil
}
