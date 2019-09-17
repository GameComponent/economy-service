package v1

import (
	"context"
	"fmt"
	"strconv"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	"go.uber.org/zap"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// GetConfig gets a config
func (s *EconomyServiceServer) GetConfig(ctx context.Context, req *v1.GetConfigRequest) (*v1.GetConfigResponse, error) {
	config, err := s.ConfigRepository.Get(ctx, req.GetKey())
	if err != nil {
		return nil, status.Error(codes.NotFound, "config not found")
	}

	return &v1.GetConfigResponse{
		Config: config,
	}, nil
}

// SetConfig sets a config
func (s *EconomyServiceServer) SetConfig(ctx context.Context, req *v1.SetConfigRequest) (*v1.SetConfigResponse, error) {
	config, err := s.ConfigRepository.Set(ctx, req.GetKey(), req.GetValue())
	if err != nil {
		s.Logger.Error("unable to set config", zap.Error(err))
		return nil, status.Error(codes.Internal, "unable to set config")
	}

	return &v1.SetConfigResponse{
		Config: config,
	}, nil
}

// ListConfig lists configs
func (s *EconomyServiceServer) ListConfig(ctx context.Context, req *v1.ListConfigRequest) (*v1.ListConfigResponse, error) {
	fmt.Println("ListConfig")

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

	// Get the items from the repository
	configs, totalSize, err := s.ConfigRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to retrieve config list")
	}

	// Determine if there is a next page
	var nextPageToken string
	if totalSize > (offset + limit) {
		nextPage := int32(parsedToken) + 1
		nextPageToken = strconv.Itoa(int(nextPage))
	}

	return &v1.ListConfigResponse{
		Configs:       configs,
		TotalSize:     totalSize,
		NextPageToken: nextPageToken,
	}, nil
}
