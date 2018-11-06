package server

import (
	"context"
	"time"

	"github.com/auburnhacks/sponsor/pkg/admin"
	"github.com/auburnhacks/sponsor/pkg/auth"
	api "github.com/auburnhacks/sponsor/proto"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// CreateAdmin is a method on the SponsorServer that is used to create an admin and save it to the database
func (ss *SponsorServer) CreateAdmin(ctx context.Context, req *api.CreateAdminRequest) (*api.CreateAdminResponse, error) {
	admin := admin.New(req.Name, req.Email, req.PasswordPlainText)
	if err := admin.Register(); err != nil {
		log.Errorf("error while registering admin: %v", err)
		return nil, err
	}
	log.Debugf("%+v", admin)
	return &api.CreateAdminResponse{
		Admin: &api.Admin{
			Name:    admin.Name,
			Email:   admin.Email,
			AdminID: admin.ID,
			ACL:     admin.ACL,
		},
	}, nil
}

// GetAdmin is a method on the SponsorServer that is used to get information of an admin
func (ss *SponsorServer) GetAdmin(ctx context.Context, req *api.GetAdminRequest) (*api.GetAdminResponse, error) {
	admin, err := admin.ByID(req.AdminID)
	if err != nil {
		return nil, err
	}
	return &api.GetAdminResponse{
		Admin: &api.Admin{
			AdminID: admin.ID,
			Name:    admin.Name,
			Email:   admin.Email,
			ACL:     admin.ACL,
		},
	}, nil
}

// DeleteAdmin is a method on the SponsorServer that is used to delete an admin from the database
func (ss *SponsorServer) DeleteAdmin(ctx context.Context, req *api.DeleteAdminRequest) (*api.DeleteAdminResponse, error) {
	return &api.DeleteAdminResponse{}, nil
}

// LoginAdmin is a methodd on the SponsorServer that is used to sign in an admin
// this method also deals with allocating a signed JWT token to the client for
// making authenticated requests
func (ss *SponsorServer) LoginAdmin(ctx context.Context, req *api.LoginAdminRequest) (*api.LoginAdminResponse, error) {
	admin, err := admin.Login(req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	issuedAt := time.Now().Unix()
	expiresAt := time.Now().AddDate(0, 0, 30).Unix()
	c := auth.NewAdminClaims(admin.ID, tokenIssuer, admin.ACL, issuedAt, expiresAt)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenStr, err := token.SignedString(ss.jwtKey)
	if err != nil {
		return nil, err
	}
	return &api.LoginAdminResponse{
		Token: tokenStr,
		Admin: &api.Admin{
			AdminID: admin.ID,
			Name:    admin.Name,
			Email:   admin.Email,
			ACL:     admin.ACL,
		},
	}, nil
}

// UpdateAdmin is a method on the SponsorServer that updates the modified state
// of an admin to the database
func (ss *SponsorServer) UpdateAdmin(ctx context.Context, req *api.UpdateAdminRequest) (*api.UpdateAdminResponse, error) {
	log.Debugf("%+v", req)
	admin, err := admin.ByID(req.AdminID)
	if err != nil {
		return nil, err
	}
	// Update all the fields
	admin.Name = req.Admin.Name
	admin.Email = req.Admin.Email
	admin.ACL = req.Admin.ACL
	if err := admin.Save(); err != nil {
		return nil, errors.Wrap(err, "error while saving admin to db")
	}
	return &api.UpdateAdminResponse{
		Admin: &api.Admin{
			AdminID: admin.ID,
			Name:    admin.Name,
			Email:   admin.Email,
			ACL:     admin.ACL,
		},
	}, nil
}