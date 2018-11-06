// Package utils provides some utility function for the entire api server
// it houses all the code for the middlewares that are deployed on the
// gRPC server
package utils

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/auburnhacks/sponsor/pkg/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

var (
	// ErrNotAuthenicated is a default error that is sent to the user if they are
	// trying to access a secure RPC call
	ErrNotAuthenicated = errors.New("auth: user is not authenticated for this RPC call")
	secret             []byte
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

func authenticateJWTToken(tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, &auth.AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		if secret == nil {
			bb, err := auth.LoadJWTKey(filepath.Join(".", "jwt_key_dev"))
			if err != nil {
				return nil, err
			}
			secret = bb
		}
		return secret, nil
	})
	if err != nil {
		return errors.Wrap(err, "auth: error while parsing token with claims")
	}
	_, ok := token.Claims.(*auth.AdminClaims)
	if !ok || !token.Valid {
		return ErrNotAuthenicated
	}
	return nil
}
