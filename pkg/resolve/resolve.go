package resolve

import (
	"context"
	"net"

	"golang.org/x/sync/singleflight"
)

type lookUp interface {
	LookupIPAddr(context.Context, string) ([]net.IPAddr, error)
}

type DefaultResolver struct {
	g singleflight.Group
	l lookUp
}

var _ Resolver = (*DefaultResolver)(nil)

func (r *DefaultResolver) Resolve(ctx context.Context, host string) ([]net.IP, error) {

	resolve := func() (interface{}, error) {
		ips, err := r.l.LookupIPAddr(ctx, host)
		if err != nil {
			return nil, err
		}

		results := make([]net.IP, len(ips))
		for i, ip := range ips {
			results[i] = ip.IP
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			return results, nil
		}

	}

	do, err, _ := r.g.Do(host, resolve)
	if err != nil {
		return nil, err
	}

	return do.([]net.IP), nil
}
