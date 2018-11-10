package server

import (
	"context"
	"fmt"

	"github.com/auburnhacks/sponsor/pkg/admin"
	"github.com/auburnhacks/sponsor/pkg/auth"
	api "github.com/auburnhacks/sponsor/proto"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// CreateAdmin is a method on the rpcServer that is used to create an admin and save it to the database
func (s *rpcServer) CreateAdmin(ctx context.Context, req *api.CreateAdminRequest) (*api.CreateAdminResponse, error) {
	a, _ := admin.ByEmail(req.Email)
	if a != nil {
		return nil, fmt.Errorf("email %s already exists", req.Email)
	}
	a = admin.New(req.Name, req.Email, req.PasswordPlainText)
	if err := a.Register(); err != nil {
		log.Errorf("error while registering admin: %v", err)
		return nil, err
	}
	log.Debugf("%+v", a)
	return &api.CreateAdminResponse{
		Admin: &api.Admin{
			Id:    a.ID,
			Name:  a.Name,
			Email: a.Email,
			ACL:   a.ACL,
		},
	}, nil
}

// GetAdmin is a method on the rpcServer that is used to get information of an admin
func (s *rpcServer) GetAdmin(ctx context.Context, req *api.GetAdminRequest) (*api.GetAdminResponse, error) {
	admin, err := admin.ByID(req.AdminId)
	if err != nil {
		return nil, err
	}
	return &api.GetAdminResponse{
		Admin: &api.Admin{
			Id:    admin.ID,
			Name:  admin.Name,
			Email: admin.Email,
			ACL:   admin.ACL,
		},
	}, nil
}

// DeleteAdmin is a method on the rpcServer that is used to delete an admin from the database
func (s *rpcServer) DeleteAdmin(ctx context.Context, req *api.DeleteAdminRequest) (*api.DeleteAdminResponse, error) {
	return &api.DeleteAdminResponse{}, nil
}

// LoginAdmin is a methodd on te rpcServer that is used to sign in an admin
// this method also deals with allocating a signed JWT token to the client for
// making authenticated requests
func (s *rpcServer) LoginAdmin(ctx context.Context, req *api.LoginAdminRequest) (*api.LoginAdminResponse, error) {
	admin, err := admin.Login(req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	cl := auth.New(admin.ID, admin.ACL)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cl)
	tokenStr, err := token.SignedString(s.privKey)
	if err != nil {
		return nil, err
	}
	return &api.LoginAdminResponse{
		Token: tokenStr,
		Admin: &api.Admin{
			Id:    admin.ID,
			Name:  admin.Name,
			Email: admin.Email,
			ACL:   admin.ACL,
		},
	}, nil
}

// UpdateAdmin is a method on te rpcServer that updates the modified state
// of an admin to the database
func (s *rpcServer) UpdateAdmin(ctx context.Context, req *api.UpdateAdminRequest) (*api.UpdateAdminResponse, error) {
	log.Debugf("%+v", req)
	admin, err := admin.ByID(req.AdminId)
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
			Id:    admin.ID,
			Name:  admin.Name,
			Email: admin.Email,
			ACL:   admin.ACL,
		},
	}, nil
}
