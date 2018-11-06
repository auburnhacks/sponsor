// Package auth provides all authentication and authorization
// functionality needed by the sponsor api server
package auth

import (
	"errors"
	"io/ioutil"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var (
	// ErrUnauthorized is an error that is used when a user is requesting an
	// unauthorized claim in the system
	ErrUnauthorized = errors.New("auth: user unauthoized to perform this action")
)

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
func HasAccessToResource(ACL string, claim string) error {
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

// NewAdminClaims returns a instance of the admin claims from the input parameters
func NewAdminClaims(adminID, issuer, acl string, issuedAt int64, expiresAt int64) *AdminClaims {
	ac := new(AdminClaims)
	ac.ACL = acl
	ac.Id = adminID
	ac.IssuedAt = issuedAt
	ac.ExpiresAt = expiresAt
	return ac
}
