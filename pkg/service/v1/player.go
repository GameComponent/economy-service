package v1

import (
	"context"
	"fmt"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
)

func (s *economyServiceServer) GetPlayer(ctx context.Context, req *v1.GetPlayerRequest) (*v1.GetPlayerResponse, error) {
	fmt.Println("GetPlayer")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	player, err := s.playerRepository.Get(req.GetPlayerId())

	if err != nil {
		return nil, err
	}

	return &v1.GetPlayerResponse{
		Api:    apiVersion,
		Player: player,
	}, nil
}
