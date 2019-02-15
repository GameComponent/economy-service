package v1

import (
	"context"
	"fmt"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	"github.com/google/uuid"
)

func (s *economyServiceServer) GiveItem(ctx context.Context, req *v1.GiveItemRequest) (*v1.GiveItemResponse, error) {
	fmt.Println("GiveItem")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	item, err := s.itemRepository.Get(ctx, req.GetItemId())
	if err != nil {
		return nil, err
	}

	storageItemID, err := s.storageRepository.GiveItem(ctx, req.GetStorageId(), req.GetItemId())
	if err != nil {
		return nil, err
	}

	storageItem := &v1.StorageItem{
		Id:   *storageItemID,
		Item: item,
		// Metadata: metadata,
	}

	return &v1.GiveItemResponse{
		Api:       apiVersion,
		StorageId: req.GetStorageId(),
		Item:      storageItem,
	}, nil
}

func (s *economyServiceServer) CreateStorage(ctx context.Context, req *v1.CreateStorageRequest) (*v1.CreateStorageResponse, error) {
	fmt.Println("CreateStorage")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	storage, err := s.storageRepository.Create(ctx, req.GetPlayerId(), req.GetName())
	if err != nil {
		return nil, err
	}

	return &v1.CreateStorageResponse{
		Api:     apiVersion,
		Storage: storage,
	}, nil
}

func (s *economyServiceServer) GetStorage(ctx context.Context, req *v1.GetStorageRequest) (*v1.GetStorageResponse, error) {
	fmt.Println("GetStorage")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// Check if the requested storage id is a valid UUID
	_, err := uuid.Parse(req.GetStorageId())
	if err != nil {
		return nil, err
	}

	storage, err := s.storageRepository.Get(ctx, req.GetStorageId())
	if err != nil {
		return nil, err
	}

	return &v1.GetStorageResponse{
		Api:     apiVersion,
		Storage: storage,
	}, nil
}

func (s *economyServiceServer) GiveCurrency(ctx context.Context, req *v1.GiveCurrencyRequest) (*v1.GiveCurrencyResponse, error) {
	fmt.Println("GiveCurrency")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	storageCurrency, err := s.storageRepository.GiveCurrency(ctx, req.GetStorageId(), req.GetCurrencyId(), req.GetAmount())
	if err != nil {
		return nil, err
	}

	return &v1.GiveCurrencyResponse{
		Api:      apiVersion,
		Currency: storageCurrency,
	}, nil
}
