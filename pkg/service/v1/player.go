package v1

import (
	"context"
	"fmt"
	"strconv"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
)

func (s *economyServiceServer) GetPlayer(ctx context.Context, req *v1.GetPlayerRequest) (*v1.GetPlayerResponse, error) {
	fmt.Println("GetPlayer")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	player, err := s.playerRepository.Get(ctx, req.GetPlayerId())

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve player")
	}

	return &v1.GetPlayerResponse{
		Api:    apiVersion,
		Player: player,
	}, nil
}

func (s *economyServiceServer) ListPlayer(ctx context.Context, req *v1.ListPlayerRequest) (*v1.ListPlayerResponse, error) {
	fmt.Println("ListPlayer")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// Parse the page token
	var parsedToken int64
	parsedToken, _ = strconv.ParseInt(req.GetPageToken(), 10, 32)

	// Get the limit
	limit := req.GetPageSize()
	if limit == 0 {
		limit = 100
	}

	// Get the offset
	offset := int32(0)
	if len(req.GetPageToken()) > 0 {
		offset = int32(parsedToken) * limit
	}

	// Get the players
	players, totalSize, err := s.playerRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Determine if there is a next page
	var nextPageToken string
	if totalSize > (offset + limit) {
		nextPage := int32(parsedToken) + 1
		nextPageToken = strconv.Itoa(int(nextPage))
	}

	return &v1.ListPlayerResponse{
		Api:           apiVersion,
		Players:       players,
		TotalSize:     totalSize,
		NextPageToken: nextPageToken,
	}, nil
}
