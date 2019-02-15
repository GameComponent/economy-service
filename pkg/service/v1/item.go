package v1

import (
	"context"
	"fmt"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
)

func (s *economyServiceServer) CreateItem(ctx context.Context, req *v1.CreateItemRequest) (*v1.CreateItemResponse, error) {
	fmt.Println("CreateItem")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// Add item to the databased return the generated UUID
	item, err := s.itemRepository.Create(ctx, req.GetName())
	if err != nil {
		return nil, err
	}

	return &v1.CreateItemResponse{
		Api:  apiVersion,
		Item: item,
	}, nil
}

func (s *economyServiceServer) UpdateItem(ctx context.Context, req *v1.UpdateItemRequest) (*v1.UpdateItemResponse, error) {
	fmt.Println("UpdateItem")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	item, err := s.itemRepository.Update(ctx, req.GetItemId(), req.GetName(), `{"kaas":"baas"}`)

	if err != nil {
		return nil, err
	}

	return &v1.UpdateItemResponse{
		Api:  apiVersion,
		Item: item,
	}, nil
}

func (s *economyServiceServer) GetItem(ctx context.Context, req *v1.GetItemRequest) (*v1.GetItemResponse, error) {
	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	item, err := s.itemRepository.Get(ctx, req.GetItemId())
	if err != nil {
		return nil, err
	}

	return &v1.GetItemResponse{
		Api:  apiVersion,
		Item: item,
	}, nil
}

func (s *economyServiceServer) ListItems(ctx context.Context, req *v1.ListItemsRequest) (*v1.ListItemsResponse, error) {
	fmt.Println("ListItems")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	items, err := s.itemRepository.List(ctx)

	if err != nil {
		return nil, err
	}

	return &v1.ListItemsResponse{
		Api:   apiVersion,
		Items: items,
	}, nil
}
