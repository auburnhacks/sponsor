package server

import (
	"context"

	"github.com/auburnhacks/sponsor/pkg/participant"
	api "github.com/auburnhacks/sponsor/proto"
)

// ListParticipants is a method on the SponsorServer that lists all the participants
// are synced from the external database
func (s *rpcServer) ListParticipants(ctx context.Context,
	req *api.ListParticipantsRequest) (*api.ListParticipantsResponse, error) {
	pSlice, err := participant.List()
	if err != nil {
		return nil, err
	}
	// prp = protobufParticipant
	prp := make([]*api.Participant, len(pSlice))
	i := 0
	for _, p := range pSlice {
		prp[i] = &api.Participant{
			Id:       p.ID,
			Name:     p.Name,
			Github:   p.Github,
			Linkedin: p.Linkedin,
			Resume:   p.Resume,
		}
		i++
	}
	return &api.ListParticipantsResponse{
		Participants: prp,
	}, nil
}
