package server

import (
	"context"

	"github.com/auburnhacks/sponsor/pkg/participant"
	api "github.com/auburnhacks/sponsor/proto"
	log "github.com/sirupsen/logrus"
)

// ListParticipants is a method on the SponsorServer that lists all the participants
// are synced from the external database
func (s *rpcServer) ListParticipants(ctx context.Context,
	req *api.ListParticipantsRequest) (*api.ListParticipantsResponse, error) {
	pSlice, err := participant.List()
	if err != nil {
		return nil, err
	}
	log.Debugf("%+v", pSlice)
	// prp = protobufParticipants
	prp := make([]*api.Participant, len(pSlice))
	i := 0
	for _, p := range pSlice {
		log.Debugf("%v", p)
		prp[i] = &api.Participant{}
	}
	return &api.ListParticipantsResponse{
		Participants: prp,
	}, nil
}
