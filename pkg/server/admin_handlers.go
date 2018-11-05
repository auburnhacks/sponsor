package server

import (
	"context"
	"time"

	"github.com/auburnhacks/sponsor/pkg/admin"
	api "github.com/auburnhacks/sponsor/proto"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

func (s *SponsorServer) CreateAdmin(ctx context.Context, req *api.CreateAdminRequest) (*api.CreateAdminResponse, error) {
	admin := admin.New(req.Name, req.Email, req.PasswordPlainText)
	if err := admin.Register(); err != nil {
		log.Errorf("error while registering admin: %v", err)
		return nil, err
	}
	log.Debugf("%+v", admin)
	return &api.CreateAdminResponse{
		Admin: &api.Admin{
			Email:   admin.Email,
			AdminID: int64(admin.ID),
			ACL:     admin.ACL,
		},
	}, nil
}

func (s *SponsorServer) GetAdmin(ctx context.Context, req *api.GetAdminRequest) (*api.GetAdminResponse, error) {
	return &api.GetAdminResponse{}, nil
}

func (s *SponsorServer) DeleteAdmin(ctx context.Context, req *api.DeleteAdminRequest) (*api.DeleteAdminResponse, error) {
	return &api.DeleteAdminResponse{}, nil
}

func (s *SponsorServer) LoginAdmin(ctx context.Context, req *api.LoginAdminRequest) (*api.LoginAdminResponse, error) {
	// TODO: look up database for the user
	admin, err := admin.Login(req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	log.Debugf("%+v", admin)
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Date(2019, 9, 30, 0, 0, 0, 0, time.UTC).Unix(),
		Issuer:    "sponsor_auburnhacks",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(signingKey)
	if err != nil {
		return nil, err
	}
	return &api.LoginAdminResponse{
		Token: tokenStr,
		Admin: &api.Admin{
			AdminID: int64(admin.ID),
			Name:    admin.Name,
			Email:   admin.Email,
			ACL:     admin.ACL,
		},
	}, nil
}
