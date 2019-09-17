package v1

import (
	"context"
	"strconv"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetPlayer gets a player
func (s *EconomyServiceServer) GetPlayer(ctx context.Context, req *v1.GetPlayerRequest) (*v1.GetPlayerResponse, error) {
	player, err := s.PlayerRepository.Get(ctx, req.GetPlayerId())

	if err != nil {
		s.Logger.Error("player not found", zap.Error(err))
		return nil, status.Error(codes.NotFound, "player not found")
	}

	return &v1.GetPlayerResponse{
		Player: player,
	}, nil
}

// CreatePlayer creates a player
func (s *EconomyServiceServer) CreatePlayer(ctx context.Context, req *v1.CreatePlayerRequest) (*v1.CreatePlayerResponse, error) {
	if req.GetPlayerId() == "" {
		return nil, status.Error(codes.InvalidArgument, "no player_id given")
	}

	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "no name given")
	}

	player, err := s.PlayerRepository.Create(
		ctx,
		req.GetPlayerId(),
		req.GetName(),
		req.GetMetadata(),
	)
	if err != nil {
		return nil, status.Error(codes.Aborted, "unable to create player, make sure the player_id is unique")
	}

	return &v1.CreatePlayerResponse{
		Player: player,
	}, nil
}

// UpdatePlayer updates a player
func (s *EconomyServiceServer) UpdatePlayer(ctx context.Context, req *v1.UpdatePlayerRequest) (*v1.UpdatePlayerResponse, error) {
	if req.GetPlayerId() == "" {
		return nil, status.Error(codes.InvalidArgument, "no player_id given")
	}

	player, err := s.PlayerRepository.Update(
		ctx,
		req.GetPlayerId(),
		req.GetName(),
		req.GetMetadata(),
	)
	if err != nil {
		s.Logger.Error("unable to update player", zap.Error(err))
		return nil, status.Error(codes.NotFound, "unable to update player")
	}

	return &v1.UpdatePlayerResponse{
		Player: player,
	}, nil
}

// ListPlayer lists players
func (s *EconomyServiceServer) ListPlayer(ctx context.Context, req *v1.ListPlayerRequest) (*v1.ListPlayerResponse, error) {
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
	players, totalSize, err := s.PlayerRepository.List(ctx, limit, offset)
	if err != nil {
		s.Logger.Error("unable to list players", zap.Error(err))
		return nil, status.Error(codes.Aborted, "unable to list players")
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

// SearchPlayer search for a player
func (s *EconomyServiceServer) SearchPlayer(ctx context.Context, req *v1.SearchPlayerRequest) (*v1.SearchPlayerResponse, error) {
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
	players, totalSize, err := s.PlayerRepository.Search(ctx, req.GetQuery(), limit, offset)
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
