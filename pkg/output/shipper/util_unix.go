//go:build !windows

package shipper

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func defaultDialOptions() []grpc.DialOption {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	return opts
}
