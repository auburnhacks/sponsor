// Package auth provides all authentication and authorization
// functionality needed by the sponsor api server
package auth

import (
	"io/ioutil"

	"github.com/dgrijalva/jwt-go"
)

// LoadJWTKey loads a JWT key from a given file path
func LoadJWTKey(filename string) ([]byte, error) {
	bb, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return bb, nil
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

// Valid is a function that is required by the jwt.Claims interface
// func (ac *AdminClaims) Valid() error {
// 	return nil
// }
