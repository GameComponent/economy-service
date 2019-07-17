package v1

import (
	"context"
	"fmt"
	"strconv"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func (s *economyServiceServer) CreateItem(ctx context.Context, req *v1.CreateItemRequest) (*v1.CreateItemResponse, error) {
	fmt.Println("CreateItem")

	// Add item to the databased return the generated UUID
	item, err := s.itemRepository.Create(
		ctx,
		req.GetName(),
		req.GetStackable(),
		req.GetStackMaxAmount(),
		int64(req.GetStackBalancingMethod()),
		req.GetMetadata(),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to create item")
	}

	return &v1.CreateItemResponse{
		Item: item,
	}, nil
}

func (s *economyServiceServer) UpdateItem(ctx context.Context, req *v1.UpdateItemRequest) (*v1.UpdateItemResponse, error) {
	fmt.Println("UpdateItem")

	item, err := s.itemRepository.Update(
		ctx,
		req.GetItemId(),
		req.GetName(),
		req.GetMetadata(),
	)

	if err != nil {
		return nil, status.Error(codes.Internal, "unable to update item")
	}

	return &v1.UpdateItemResponse{
		Item: item,
	}, nil
}

func (s *economyServiceServer) GetItem(ctx context.Context, req *v1.GetItemRequest) (*v1.GetItemResponse, error) {
	item, err := s.itemRepository.Get(ctx, req.GetItemId())
	if err != nil {
		return nil, err
	}

	return &v1.GetItemResponse{
		Item: item,
	}, nil
}

func (s *economyServiceServer) ListItem(ctx context.Context, req *v1.ListItemRequest) (*v1.ListItemResponse, error) {
	fmt.Println("ListItem")

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
	items, totalSize, err := s.itemRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to retrieve item list")
	}

	// Determine if there is a next page
	var nextPageToken string
	if totalSize > (offset + limit) {
		nextPage := int32(parsedToken) + 1
		nextPageToken = strconv.Itoa(int(nextPage))
	}

	return &v1.ListItemResponse{
		Items:         items,
		TotalSize:     totalSize,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *economyServiceServer) SearchItem(ctx context.Context, req *v1.SearchItemRequest) (*v1.SearchItemResponse, error) {
	fmt.Println("SearchItem")

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

	// Search the items
	items, totalSize, err := s.itemRepository.Search(ctx, req.GetQuery(), limit, offset)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to retrieve item search results")
	}

	// Determine if there is a next page
	var nextPageToken string
	if totalSize > (offset + limit) {
		nextPage := int32(parsedToken) + 1
		nextPageToken = strconv.Itoa(int(nextPage))
	}

	return &v1.SearchItemResponse{
		Items:         items,
		TotalSize:     totalSize,
		NextPageToken: nextPageToken,
	}, nil
}
