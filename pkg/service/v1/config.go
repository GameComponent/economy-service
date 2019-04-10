package v1

import (
	"context"
	"fmt"

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
