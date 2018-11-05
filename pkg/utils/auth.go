// Package utils provides some utility function for the entire api server
// it houses all the code for the middlewares that are deployed on the
// gRPC server
package utils

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/metadata"
)

var (
	// ErrNotAuthenicated is a default error that is sent to the user if they are
	// trying to access a secure RPC call
	ErrNotAuthenicated = errors.New("user is not autenticated for this RPC call")
)

func authenticate(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("error gathering metadata from context")
	}
	token, ok := md["authorization"]
	if !ok {
		return ErrNotAuthenicated
	}
	bearerToken := strings.Split(token[0], " ")
	return authenticateJWTToken(bearerToken[1])
}

func authenticateJWTToken(token string) error {
	// TODO: add all the JWT authentication stuff here
	return nil
}
