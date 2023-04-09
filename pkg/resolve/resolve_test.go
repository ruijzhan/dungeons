package resolve

import (
	"context"
	"errors"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/ruijzhan/dungeons/pkg/cache"
)

type fakeCache struct{}

func (fakeCache) Get(key any) (any, error) {
	return nil, errors.New("fake not implemented")
}

func (fakeCache) Set(key any, value any) error {
	return errors.New("fake not implemented")
}

func TestResolver_Resolve(t *testing.T) {
	resolver := &DefaultResolver{
		lookup: net.DefaultResolver,
		cache:  fakeCache{},
	}

	t.Run("success", func(t *testing.T) {
		ips, err := resolver.Resolve(context.Background(), "example.com")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(ips) == 0 {
			t.Errorf("expected non-empty slice of IP addresses, but got empty slice: %v", ips)
		}
		for _, ip := range ips {
			if ip.To4() == nil && ip.To16() == nil {
				t.Errorf("expected valid IP address, but got invalid address: %s", ip)
			}
		}
	})

	t.Run("canceled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		ips, err := resolver.Resolve(ctx, "example.com")
		if err == nil {
			t.Errorf("expected error, but got none")
		}
		if ips != nil {
			t.Errorf("expected nil slice of IP addresses, but got: %v", ips)
		}
	})

	t.Run("invalid host", func(t *testing.T) {
		ips, err := resolver.Resolve(context.Background(), "not-a-real-hostname")
		if err == nil {
			t.Errorf("expected error, but got none")
		}
		if ips != nil {
			t.Errorf("expected nil slice of IP addresses, but got: %v", ips)
		}
	})
}

type mockResolver struct {
	ip net.IP
}

func (m *mockResolver) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	return []net.IPAddr{{IP: m.ip}}, nil
}

func TestDefaultResolver_Resolve(t *testing.T) {
	// Create a test context with a short timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	tests := []struct {
		name     string
		host     string
		ip       net.IP
		timeout  bool
		err      error
		resolver lookUp
	}{
		{
			name: "successful lookup",
			host: "foo.com",
			ip:   net.ParseIP("192.0.2.1"),
			// Replace the default resolver with a mock that returns the desired IP address.
			resolver: &mockResolver{ip: net.ParseIP("192.0.2.1")},
		},
		{
			name:     "failed lookup",
			host:     "barxxxx.comxxx",
			err:      errors.New("lookup barxxxx.comxxx: no such host"),
			resolver: net.DefaultResolver,
		},
		{
			name:     "context deadline exceeded",
			host:     "baz.com",
			timeout:  true,
			err:      errors.New("operation was canceled"),
			resolver: net.DefaultResolver,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := DefaultResolver{
				lookup: tt.resolver,
				cache:  fakeCache{},
			}

			if tt.timeout {
				cancel()
			}

			got, err := r.Resolve(ctx, tt.host)

			if tt.err != nil && err != nil {
				if !strings.Contains(err.Error(), tt.err.Error()) {
					t.Errorf("expected error %q, but got %q", tt.err, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(got) != 1 || !got[0].Equal(tt.ip) {
				t.Errorf("expected IP %v, but got %v", tt.ip, got)
			}
		})
	}
}

type mockLookup struct{}

func (m mockLookup) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	switch host {
	case "example.com":
		return []net.IPAddr{{IP: net.ParseIP("93.184.216.34")}}, nil
	case "localhost":
		return []net.IPAddr{{IP: net.ParseIP("127.0.0.1")}, {IP: net.ParseIP("::1")}}, nil
	default:
		return nil, &net.DNSError{Err: "lookup failed", Name: host}
	}
}

func TestDefaultResolver_Resolve2(t *testing.T) {
	ctx := context.Background()
	r := DefaultResolver{
		lookup: &mockLookup{},
		cache:  cache.NewTTL(),
	}

	t.Run("cache hit", func(t *testing.T) {
		ip := net.ParseIP("93.184.216.34")
		r.cache.Set("example.com", []net.IP{ip})
		res, err := r.Resolve(ctx, "example.com")
		if err != nil {
			t.Errorf("Expected no error but got %v", err)
		}
		if len(res) != 1 || !res[0].Equal(ip) {
			t.Errorf("Expected %v but got %v", []net.IP{ip}, res)
		}
	})

	t.Run("cache miss", func(t *testing.T) {
		res, err := r.Resolve(ctx, "example.com")
		if err != nil {
			t.Errorf("Expected no error but got %v", err)
		}
		if len(res) != 1 || !res[0].Equal(net.ParseIP("93.184.216.34")) {
			t.Errorf("Expected [93.184.216.34] but got %v", res)
		}
	})

	t.Run("lookup failed", func(t *testing.T) {
		_, err := r.Resolve(ctx, "nonexistent-host")
		if err == nil {
			t.Errorf("Expected an error but got none")
		}
	})
}
