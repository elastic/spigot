//go:build windows

package shipper

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/elastic/elastic-agent-libs/api/npipe"
)

func defaultDialOptions() []grpc.DialOption {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(npipeDialer),
	}
	return opts
}

func npipeDialer(ctx context.Context, addr string) (net.Conn, error) {
	return npipe.DialContext(addr)(ctx, "", "")
}
