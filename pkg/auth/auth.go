// Package auth provides all authentication and authorization
// functionality needed by the sponsor api server
package auth

import (
	"context"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/metadata"
)

var (
	// ErrUnauthorized is an error that is used when a user is requesting an
	// unauthorized claim in the system
	ErrUnauthorized = errors.New("auth: user unauthoized to perform this action")
)

// Claims is an interface that describes the various cliams that are
// allowed in the auth package
type Claims interface {
	Claim(resource string) error
	jwt.Claims
}

// LoadJWTKey loads a JWT key from a given file path
func LoadJWTKey(filename string) ([]byte, error) {
	bb, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return bb, nil
}

// HasAccessToResource is a function that determines whether an
// ACL has access to the requesting claim
// NOTE: ACL is has to be a comma separated string
func hasAccessToResource(ACL string, claim string) error {
	acls := strings.Split(ACL, ",")
	for _, acl := range acls {
		if claim == acl {
			return nil
		}
	}
	return ErrUnauthorized
}

// AdminClaims is a struct that represents the structure
// of the JWT token cliams. It implements Claims interface
// from the JWT package
type AdminClaims struct {
	ACL string `json:"acl"`
	jwt.StandardClaims
}

// Claim is a function that returns an error if the requested
// resource is not present in the acl
func (c *AdminClaims) Claim(resource string) error {
	return hasAccessToResource(c.ACL, resource)
}

// newAdminClaims returns a instance of the admin claims from the input parameters
func newAdminClaims(adminID, issuer, acl string, issuedAt int64, expiresAt int64) *AdminClaims {
	ac := new(AdminClaims)
	ac.ACL = acl
	ac.Id = adminID
	ac.IssuedAt = issuedAt
	ac.ExpiresAt = expiresAt
	return ac
}

// New returns a struct that implements the Claims interface
func New(id, acl string) Claims {
	issuedAt := time.Now().Unix()
	expiresAt := time.Now().AddDate(0, 0, 30).Unix()
	return newAdminClaims(id, "sponsor_auburnhacks", acl, issuedAt, expiresAt)
}

// FromContext returns the cliams madde by the user based on the request context
func FromContext(ctx context.Context) (Claims, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("auth: error while gathering metadata")
	}
	authStr, ok := md["authorization"]
	if !ok {
		return nil, errors.New("auth: not authorization found in request context")
	}
	beaerToken := strings.Split(authStr[0], " ")
	token, err := jwt.ParseWithClaims(beaerToken[1], &AdminClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return LoadJWTKey(filepath.Join(".", "jwt_key_dev"))
		},
	)
	if err != nil {
		return nil, err
	}
	cl, ok := token.Claims.(*AdminClaims)
	if !ok {
		return nil, errors.New("auth: error converting to admin claims")
	}
	return cl, nil
}
