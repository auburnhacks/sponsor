package server

import (
	"context"

	api "github.com/auburnhacks/sponsor/proto"
)

// ListParticipants is a method on the SponsorServer that lists all the participants
// are synced from the external database
func (ss *SponsorServer) ListParticipants(ctx context.Context, req *api.ListParticipantsRequest) (*api.ListParticipantsResponse, error) {

	return &api.ListParticipantsResponse{}, nil
}
