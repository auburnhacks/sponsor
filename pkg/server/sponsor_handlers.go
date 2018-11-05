package server

import (
	"context"

	"github.com/auburnhacks/sponsor/pkg/sponsor"
	api "github.com/auburnhacks/sponsor/proto"
)

// CreateSponsor is a method on the SponsorServer that is used to create a sponsor
// this is typically called by an admin
// NOTE (kirandasika98): This can change later based on a new feature change
func (ss *SponsorServer) CreateSponsor(ctx context.Context,
	req *api.CreateSponsorRequest) (*api.CreateSponsorResponse, error) {
	c, err := sponsor.CompanyByID(req.Sponsor.Company.Id)
	if err != nil {
		return nil, err
	}
	s := sponsor.New(req.Sponsor.Name, req.Sponsor.Email, req.Sponsor.Password, c, req.Sponsor.ACL)
	if err := s.Register(); err != nil {
		return nil, err
	}
	return &api.CreateSponsorResponse{
		Sponsor: &api.Sponsor{
			Id:    s.ID,
			Name:  s.Name,
			Email: s.Email,
			Company: &api.Company{
				Id:   s.Company.Name,
				Name: s.Company.Name,
				Logo: s.Company.Logo,
			},
			ACL: s.ACL,
		},
	}, nil
}

// CreateCompany creates a company and saved it to the database
func (ss *SponsorServer) CreateCompany(ctx context.Context,
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
