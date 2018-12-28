package server

import (
	"context"

	"github.com/auburnhacks/sponsor/pkg/participant"
	"github.com/auburnhacks/sponsor/pkg/sponsor"
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
			Id:         p.ID,
			Name:       p.Name,
			Email:      p.Email,
			University: p.University,
			Major:      p.Major,
			GradYear:   int32(p.GradYear),
			Github:     p.Github,
			Linkedin:   p.Linkedin,
			Resume:     p.Resume,
		}
		i++
	}
	return &api.ListParticipantsResponse{
		Participants: prp,
	}, nil
}

func (s *rpcServer) ListCompanies(ctx context.Context,
	req *api.ListCompaniesRequest) (*api.ListCompaniesResponse, error) {
	companies, err := sponsor.ListCompanies()
	if err != nil {
		return nil, err
	}
	var apiCompanies []*api.Company
	for _, c := range companies {
		apiC := &api.Company{
			Id:   c.ID,
			Name: c.Name,
			Logo: c.Logo,
		}
		apiCompanies = append(apiCompanies, apiC)
	}
	return &api.ListCompaniesResponse{
		Companies: apiCompanies,
	}, nil
}
