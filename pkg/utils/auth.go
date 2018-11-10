// Package utils provides some utility function for the entire api server
// it houses all the code for the middlewares that are deployed on the
// gRPC server
package utils

import (
	"context"

	"github.com/auburnhacks/sponsor/pkg/auth"
	"github.com/pkg/errors"
)

var (
	// ErrNotAuthenicated is a default error that is sent to the user if they are
	// trying to access a secure RPC call
	ErrNotAuthenicated = errors.New("auth: user is not authenticated for this RPC call")
	secret             []byte
)

// authenticate is a helper function that the invoked by the
// gRPC interceptors to authenticate RPC requests
func authenticate(ctx context.Context) error {
	_, err := auth.FromContext(ctx)
	if err != nil {
		return err
	}
	return nil
}
