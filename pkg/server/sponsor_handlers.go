package server

import (
	"context"

	"github.com/auburnhacks/sponsor/pkg/auth"
	"github.com/auburnhacks/sponsor/pkg/participant"
	"github.com/auburnhacks/sponsor/pkg/sponsor"
	api "github.com/auburnhacks/sponsor/proto"
	jwt "github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

// LoginSponsor is a method on the rpcServer that is used to validate and
// login a sponsor and issue a JWT token
func (ss *rpcServer) LoginSponsor(ctx context.Context,
	req *api.LoginSponsorRequest) (*api.LoginSponsorResponse, error) {
	sp, err := sponsor.Login(req.Email, req.PasswordPlainText)
	if err != nil {
		return nil, err
	}
	c, err := sponsor.CompanyByID(sp.CompanyID)
	if err != nil {
		return nil, err
	}
	cl := auth.New(sp.ID, sp.ACL)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, cl)
	tokenStr, err := token.SignedString(ss.privKey)
	if err != nil {
		return nil, err
	}
	return &api.LoginSponsorResponse{
		Token: tokenStr,
		Sponsor: &api.Sponsor{
			Id:    sp.ID,
			Name:  sp.Name,
			Email: sp.Email,
			ACL:   sp.ACL,
			Company: &api.Company{
				Id:   c.ID,
				Name: c.Name,
				Logo: c.Logo,
			},
		},
	}, nil
}

// CreateSponsor is a method on the rpcServer that is used to create a sponsor
// this is typically called by an admin
// NOTE (kirandasika98): This can change later based on a new feature change
func (ss *rpcServer) CreateSponsor(ctx context.Context,
	req *api.CreateSponsorRequest) (*api.CreateSponsorResponse, error) {
	c, err := sponsor.CompanyByID(req.Sponsor.Company.Id)
	if err != nil {
		return nil, err
	}
	s := sponsor.New(req.Sponsor.Name, req.Sponsor.Email, req.Sponsor.Password,
		req.Sponsor.Company.Id, req.Sponsor.ACL)
	if err := s.Register(); err != nil {
		return nil, err
	}
	return &api.CreateSponsorResponse{
		Sponsor: &api.Sponsor{
			Id:    s.ID,
			Name:  s.Name,
			Email: s.Email,
			Company: &api.Company{
				Id:   c.ID,
				Name: c.Name,
				Logo: c.Logo,
			},
			ACL: s.ACL,
		},
	}, nil
}

// GetSponsor is a method on the rpcServer that is used to get information
// about a sponsor.
// NOTE: password is not sent through this rpc call
func (ss *rpcServer) GetSponsor(ctx context.Context,
	req *api.GetSponsorRequest) (*api.GetSponsorResponse, error) {
	s, err := sponsor.ByID(req.SponsorId)
	if err != nil {
		return nil, err
	}
	c, err := sponsor.CompanyByID(s.CompanyID)
	if err != nil {
		return nil, err
	}
	return &api.GetSponsorResponse{
		Sponsor: &api.Sponsor{
			Id:    s.ID,
			Name:  s.Name,
			Email: s.Email,
			Password: "",
			Company: &api.Company{
				Id:   c.ID,
				Name: c.Name,
				Logo: c.Logo,
			},
			ACL: s.ACL,
		},
	}, nil
}

// UpdateSponsor is a method on the rpcServer that is used to modify a
// state of a sponsor in the database
func (ss *rpcServer) UpdateSponsor(ctx context.Context, req *api.UpdateSponsorRequest) (*api.UpdateSponsorResponse, error) {
	s, err := sponsor.ByID(req.SponsorId)
	if err != nil {
		return nil, err
	}
	log.Debugf("%+v", s)
	s.Name = req.Sponsor.Name
	s.Email = req.Sponsor.Email
	s.ACL = req.Sponsor.ACL
	if err := s.Save(); err != nil {
		return nil, err
	}
	c, err := sponsor.CompanyByID(s.CompanyID)
	if err != nil {
		return nil, err
	}
	return &api.UpdateSponsorResponse{
		Sponsor: &api.Sponsor{
			Id:    s.ID,
			Name:  s.Name,
			Email: s.Email,
			Company: &api.Company{
				Id:   c.ID,
				Name: c.Name,
				Logo: c.Logo,
			},
		},
	}, nil
}

// CreateCompany creates a company and saved it to the database
func (ss *rpcServer) CreateCompany(ctx context.Context,
	req *api.CreateCompanyRequest) (*api.CreateCompanyResponse, error) {
	c := sponsor.NewCompany(req.Name, req.Logo)

	// Authorization stuff for admin
	cl, err := auth.FromContext(ctx)
	if err != nil {
		return nil, err
	}
	if err := cl.Claim("update"); err != nil {
		return nil, err
	}

	if err := c.Save(); err != nil {
		return nil, err
	}
	return &api.CreateCompanyResponse{
		Company: &api.Company{
			Id:   c.ID,
			Name: c.Name,
			Logo: c.Logo,
		},
	}, nil
}

// Resumes is a function that returns a tar archive as a sequence of
// bytes
func (ss *rpcServer) Resumes(ctx context.Context,
	req *api.ResumesRequest) (*api.ResumesResponse, error) {
	b, err := participant.AllResumes()
	if err != nil {
		return nil, err
	}
	return &api.ResumesResponse{
		Archive: b,
	}, nil
}
