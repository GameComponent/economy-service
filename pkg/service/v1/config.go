package v1

import (
	"context"
	"fmt"
	"strconv"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
)

func (s *economyServiceServer) GetConfig(ctx context.Context, req *v1.GetConfigRequest) (*v1.GetConfigResponse, error) {
	fmt.Println("GetConfig")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	config, err := s.configRepository.Get(ctx, req.GetKey())
	if err != nil {
		return nil, fmt.Errorf("key not found")
	}

	return &v1.GetConfigResponse{
		Api:    apiVersion,
		Config: config,
	}, nil
}

func (s *economyServiceServer) SetConfig(ctx context.Context, req *v1.SetConfigRequest) (*v1.SetConfigResponse, error) {
	fmt.Println("SetConfig")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	config, err := s.configRepository.Set(ctx, req.GetKey(), req.GetValue())
	if err != nil {
		return nil, err
	}

	return &v1.SetConfigResponse{
		Api:    apiVersion,
		Config: config,
	}, nil
}

func (s *economyServiceServer) ListConfig(ctx context.Context, req *v1.ListConfigRequest) (*v1.ListConfigResponse, error) {
	fmt.Println("ListConfig")

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

	// Get the items from the repository
	configs, totalSize, err := s.configRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Determine if there is a next page
	var nextPageToken string
	if totalSize > (offset + limit) {
		nextPage := int32(parsedToken) + 1
		nextPageToken = strconv.Itoa(int(nextPage))
	}

	return &v1.ListConfigResponse{
		Api:           apiVersion,
		Configs:       configs,
		TotalSize:     totalSize,
		NextPageToken: nextPageToken,
	}, nil
}
