package utils

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/metadata"
)

var (
	NotAuthenicated = errors.New("user is not autenticated for this RPC call")
)

func authenticate(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("error gathering metadata from context")
	}
	token, ok := md["authorization"]
	if !ok {
		return NotAuthenicated
	}
	bearerToken := strings.Split(token[0], " ")
	return authenticateJWTToken(bearerToken[1])
}

func authenticateJWTToken(token string) error {
	// TODO: add all the JWT authentication stuff here
	return nil
}
