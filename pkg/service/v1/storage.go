package v1

import (
	"context"
	"fmt"
	"strconv"

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

	fmt.Println("item", item)

	if item.Stackable && item.StackBalancingMethod == v1.Item_UNBALANCED_FILL_EXISTING_STACKS {
		return nil, fmt.Errorf("unimplemented")
	}

	if item.Stackable &&
		item.StackBalancingMethod != v1.Item_DEFAULT &&
		item.StackBalancingMethod != v1.Item_UNBALANCED_CREATE_NEW_STACKS {
		storage, err := s.storageRepository.Get(ctx, req.GetStorageId())
		if err != nil {
			return nil, err
		}

		existingStorageItems := []*v1.StorageItem{}
		for _, storageItem := range storage.Items {
			if storageItem.Item.Id != req.GetItemId() {
				continue
			}

			if storageItem.Item.StackMaxCount > 0 && storageItem.Amount >= storageItem.Item.StackMaxCount {
				continue
			}

			existingStorageItems = append(existingStorageItems, storageItem)
		}

		if len(existingStorageItems) > 0 {
			storageItem := existingStorageItems[0]

			err := s.storageRepository.IncreaseItemAmount(
				ctx,
				storageItem.Id,
				req.GetAmount(),
			)
			if err != nil {
				return nil, err
			}

			return &v1.GiveItemResponse{
				Api:       apiVersion,
				StorageId: req.GetStorageId(),
				Item:      storageItem,
			}, nil
		}
	}

	// For new stacks and full stack for unstackable items
	// and DEFAULT & UNBALANCED_CREATE_NEW_STACKS stack balancing methods
	storageItemID, err := s.storageRepository.GiveItem(
		ctx,
		req.GetStorageId(),
		req.GetItemId(),
		req.GetAmount(),
	)
	if err != nil {
		return nil, err
	}

	storageItem := &v1.StorageItem{
		Id:   *storageItemID,
		Item: item,
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

	if req.GetName() == "" {
		return nil, fmt.Errorf("name should not be empty")
	}

	storage, err := s.storageRepository.Create(ctx, req.GetPlayerId(), req.GetName())
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
		return nil, fmt.Errorf("unable to create storage, make sure the player exists")
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

	// Check if the request
	if req.GetStorageId() == "" {
		return nil, fmt.Errorf("the request should contain the storage_id")
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

func (s *economyServiceServer) ListStorage(ctx context.Context, req *v1.ListStorageRequest) (*v1.ListStorageResponse, error) {
	fmt.Println("ListStorage")

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
	storages, totalSize, err := s.storageRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Determine if there is a next page
	var nextPageToken string
	if totalSize > (offset + limit) {
		nextPage := int32(parsedToken) + 1
		nextPageToken = strconv.Itoa(int(nextPage))
	}

	return &v1.ListStorageResponse{
		Api:           apiVersion,
		Storages:      storages,
		NextPageToken: nextPageToken,
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
