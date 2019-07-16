package v1

import (
	"context"
	"fmt"
	"strconv"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *economyServiceServer) GetPlayer(ctx context.Context, req *v1.GetPlayerRequest) (*v1.GetPlayerResponse, error) {
	fmt.Println("GetPlayer")

	player, err := s.playerRepository.Get(ctx, req.GetPlayerId())

	if err != nil {
		return nil, status.Error(codes.NotFound, "player not found")
	}

	return &v1.GetPlayerResponse{
		Player: player,
	}, nil
}

func (s *economyServiceServer) CreatePlayer(ctx context.Context, req *v1.CreatePlayerRequest) (*v1.CreatePlayerResponse, error) {
	fmt.Println("CreatePlayer")

	if req.GetPlayerId() == "" {
		return nil, status.Error(codes.InvalidArgument, "no player_id given")
	}

	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "no name given")
	}

	player, err := s.playerRepository.Create(
		ctx,
		req.GetPlayerId(),
		req.GetName(),
	)

	if err != nil {
		return nil, status.Error(codes.Aborted, "unable to create player, make sure the player_id is unique")
	}

	return &v1.CreatePlayerResponse{
		Player: player,
	}, nil
}

func (s *economyServiceServer) UpdatePlayer(ctx context.Context, req *v1.UpdatePlayerRequest) (*v1.UpdatePlayerResponse, error) {
	fmt.Println("UpdatePlayer")

	if req.GetPlayerId() == "" {
		return nil, status.Error(codes.InvalidArgument, "no player_id given")
	}

	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "no name given")
	}

	player, err := s.playerRepository.Update(
		ctx,
		req.GetPlayerId(),
		req.GetName(),
	)
	if err != nil {
		return nil, status.Error(codes.NotFound, "player not found")
	}

	return &v1.UpdatePlayerResponse{
		Player: player,
	}, nil
}

func (s *economyServiceServer) ListPlayer(ctx context.Context, req *v1.ListPlayerRequest) (*v1.ListPlayerResponse, error) {
	fmt.Println("ListPlayer")

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
		Players:       players,
		TotalSize:     totalSize,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *economyServiceServer) SearchPlayer(ctx context.Context, req *v1.SearchPlayerRequest) (*v1.SearchPlayerResponse, error) {
	fmt.Println("SearchPlayer")

	// Check if query is empty
	if req.GetQuery() == "" {
		return nil, status.Error(codes.InvalidArgument, "no query given")
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

	// Search the players
	players, totalSize, err := s.playerRepository.Search(ctx, req.GetQuery(), limit, offset)
	if err != nil {
		return nil, err
	}

	// Determine if there is a next page
	var nextPageToken string
	if totalSize > (offset + limit) {
		nextPage := int32(parsedToken) + 1
		nextPageToken = strconv.Itoa(int(nextPage))
	}

	return &v1.SearchPlayerResponse{
		Players:       players,
		TotalSize:     totalSize,
		NextPageToken: nextPageToken,
	}, nil
}
