package server

import (
	"context"

	"github.com/auburnhacks/sponsor/pkg/sponsor"
	api "github.com/auburnhacks/sponsor/proto"
	"github.com/pkg/errors"
)

// CreateSponsor is a method on the SponsorServer that is used to create a sponsor
// this is typically called by an admin
// NOTE (kirandasika98): This can change later based on a new feature change
func (ss *SponsorServer) CreateSponsor(ctx context.Context,
	req *api.CreateSponsorRequest) (*api.CreateSponsorResponse, error) {

	s := sponsor.New(req.Sponsor.Name, req.Sponsor.Email, req.Sponsor.Password,
		req.Sponsor.Company, req.Sponsor.ACL)
	if err := s.Register(); err != nil {
		return nil, errors.Wrap(err, "error while creating sponsor")
	}
	return &api.CreateSponsorResponse{}, nil
}
