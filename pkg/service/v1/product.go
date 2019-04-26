package v1

import (
	"context"
	"fmt"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
)

func (s *economyServiceServer) CreateProduct(ctx context.Context, req *v1.CreateProductRequest) (*v1.CreateProductResponse, error) {
	fmt.Println("CreateProduct")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// Add product to the databased return the generated UUID
	product, err := s.productRepository.Create(ctx, req.GetName())
	if err != nil {
		return nil, err
	}

	return &v1.CreateProductResponse{
		Api:     apiVersion,
		Product: product,
	}, nil
}

func (s *economyServiceServer) GetProduct(ctx context.Context, req *v1.GetProductRequest) (*v1.GetProductResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	product, err := s.productRepository.Get(ctx, req.GetProductId())
	if err != nil {
		return nil, err
	}

	return &v1.GetProductResponse{
		Api:     apiVersion,
		Product: product,
	}, nil
}

func (s *economyServiceServer) AttachItem(ctx context.Context, req *v1.AttachItemRequest) (*v1.AttachItemResponse, error) {
	fmt.Println("AttachItem")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// Add product to the databased return the generated UUID
	product, err := s.productRepository.AttachItem(
		ctx,
		req.GetProductId(),
		req.GetItemId(),
		req.GetAmount(),
	)

	if err != nil {
		return nil, err
	}

	return &v1.AttachItemResponse{
		Api:     apiVersion,
		Product: product,
	}, nil
}

func (s *economyServiceServer) DetachItem(ctx context.Context, req *v1.DetachItemRequest) (*v1.DetachItemResponse, error) {
	fmt.Println("DetachItem")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// Add product to the databased return the generated UUID
	product, err := s.productRepository.DetachItem(
		ctx,
		req.GetProductItemId(),
	)

	if err != nil {
		return nil, err
	}

	return &v1.DetachItemResponse{
		Api:     apiVersion,
		Product: product,
	}, nil
}
