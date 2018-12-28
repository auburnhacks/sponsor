package utils

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/auburnhacks/sponsor/pkg/log"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// unAuthenticatedRPC is a maps all the service name to the RPC calls that
// dont' require any kind of authentication
var unauthenticatedRPC = map[string][]string{
	"proto.SponsorService": []string{
		"LoginAdmin",
		"LoginSponsor",
		"CreateAdmin",
		"ListCompanines",
	},
}

// UnaryAuthInterceptor is a gRPC middleware that intercepts all
// unary RPC calls and check whether they are authenticated
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
	// update the context with the logger
	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	reqIDHeader := metadata.Pairs("X-Request-ID", uuid.String())
	p, ok := peer.FromContext(ctx)
	if !ok {
		p = new(peer.Peer)
		p.Addr = nil
	}
	fields := logrus.Fields{
		"request_id": uuid.String(),
		"start_time": start,
		"origin_ip":  p.Addr,
	}
	ctx = log.WithFields(ctx, fields)
	grpc.SetHeader(ctx, reqIDHeader)
	logger := log.GetLogger(ctx)

	h, err := handler(ctx, req)
	if err != nil {
		logger.Errorf("%s took %d", info.FullMethod, time.Since(start))
	} else {
		logger.Infof("%s took %d", info.FullMethod, time.Since(start))
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
