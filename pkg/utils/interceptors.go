package utils

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"google.golang.org/grpc"
)

// unAuthenticatedRPC is a maps all the service name to the RPC calls that
// dont' require any kind of authentication
var unauthenticatedRPC = map[string][]string{
	"proto.SponsorService": []string{
		"LoginAdmin",
	},
}

func UnaryAuthInterceptor(ctx context.Context,
	req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// no authorization required for LoginAdmin RPC call
	if err := isUnauthenticatedRPC(info.FullMethod); err != nil {
		if err := authenticate(ctx); err != nil {
			return nil, err
		}
	}

	h, err := handler(ctx, req)
	if err != nil {
		log.Printf("request - Method:%s\tDuration:%s\tError:%v\n",
			info.FullMethod,
			time.Since(start),
			err)
	} else {
		log.Printf("request - Method:%s\tDuration:%s\n",
			info.FullMethod,
			time.Since(start))
	}
	return h, err
}

func isUnauthenticatedRPC(fullMethod string) error {
	fullMethodSlice := strings.Split(fullMethod, "/")
	rpcs, ok := unauthenticatedRPC[fullMethodSlice[1]]
	if !ok {
		return fmt.Errorf("could not service with name: %v", fullMethodSlice[1])
	}
	for _, rpc := range rpcs {
		if fullMethodSlice[2] == rpc {
			return nil
		}
	}
	return fmt.Errorf("rpc call %v not an unauthenticated call", fullMethodSlice[2])
}
