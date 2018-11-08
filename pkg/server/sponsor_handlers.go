package server

import (
	"context"

	"github.com/auburnhacks/sponsor/pkg/sponsor"
	api "github.com/auburnhacks/sponsor/proto"
	log "github.com/sirupsen/logrus"
)

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
