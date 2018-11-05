package server

import (
	"context"

	api "github.com/auburnhacks/sponsor/proto"
)

// CreateSponsor is a method on the SponsorServer that is used to create a sponsor
// this is typically called by an admin
// NOTE (kirandasika98): This can change later based on a new feature change
func (ss *SponsorServer) CreateSponsor(ctx context.Context,
	req *api.CreateSponsorRequest) (*api.CreateSponsorResponse, error) {
	return &api.CreateSponsorResponse{}, nil
}
