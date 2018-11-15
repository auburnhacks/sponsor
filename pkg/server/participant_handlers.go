package server

import (
	"context"

	api "github.com/auburnhacks/sponsor/proto"
)

// ListParticipants is a method on the SponsorServer that lists all the participants
// are synced from the external database
func (s *rpcServer) ListParticipants(ctx context.Context, req *api.ListParticipantsRequest) (*api.ListParticipantsResponse, error) {
	// pSlice, err := participant.List()
	// if err != nil {
	// 	return nil, err
	// }
	// pbParticipants := []*api.Participant{}
	return &api.ListParticipantsResponse{}, nil
}
