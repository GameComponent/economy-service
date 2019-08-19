package v1

import (
	"context"
	"fmt"
	"math"
	"strconv"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	"github.com/GameComponent/economy-service/pkg/helper/random"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *economyServiceServer) CreateStorage(ctx context.Context, req *v1.CreateStorageRequest) (*v1.CreateStorageResponse, error) {
	fmt.Println("CreateStorage")

	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "no name given")
	}

	storage, err := s.storageRepository.Create(
		ctx,
		req.GetPlayerId(),
		req.GetName(),
		req.GetMetadata(),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to create storage")
	}

	return &v1.CreateStorageResponse{
		Storage: storage,
	}, nil
}

func (s *economyServiceServer) UpdateStorage(ctx context.Context, req *v1.UpdateStorageRequest) (*v1.UpdateStorageResponse, error) {
	fmt.Println("UpdateStorage")

	storage, err := s.storageRepository.Update(
		ctx,
		req.GetStorageId(),
		req.GetName(),
		req.GetMetadata(),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to update storage")
	}

	return &v1.UpdateStorageResponse{
		Storage: storage,
	}, nil
}

func (s *economyServiceServer) GetStorage(ctx context.Context, req *v1.GetStorageRequest) (*v1.GetStorageResponse, error) {
	fmt.Println("GetStorage")

	// Check if the request
	if req.GetStorageId() == "" {
		return nil, status.Error(codes.InvalidArgument, "no storage_id given")
	}

	storage, err := s.storageRepository.Get(ctx, req.GetStorageId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "storage not found")
	}

	return &v1.GetStorageResponse{
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
		return nil, status.Error(codes.Internal, "unable to retrieve storage list")
	}

	// Determine if there is a next page
	var nextPageToken string
	if totalSize > (offset + limit) {
		nextPage := int32(parsedToken) + 1
		nextPageToken = strconv.Itoa(int(nextPage))
	}

	return &v1.ListStorageResponse{
		Storages:      storages,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *economyServiceServer) GiveCurrency(ctx context.Context, req *v1.GiveCurrencyRequest) (*v1.GiveCurrencyResponse, error) {
	fmt.Println("GiveCurrency")

	amount := random.GenerateRandomInt(
		req.GetAmount().MinAmount,
		req.GetAmount().MaxAmount,
	)

	storageCurrency, err := s.storageRepository.GiveCurrency(
		ctx,
		req.GetStorageId(),
		req.GetCurrencyId(),
		amount,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to give currency to storage")
	}

	return &v1.GiveCurrencyResponse{
		Currency: storageCurrency,
	}, nil
}

func (s *economyServiceServer) SplitStack(ctx context.Context, req *v1.SplitStackRequest) (*v1.SplitStackResponse, error) {
	fmt.Println("SplitStack")

	// Find the selected StorageItem
	selectedStorageItem := &v1.StorageItem{}
	storage, err := s.storageRepository.Get(ctx, req.GetStorageId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "storage not found")
	}

	for _, storageItem := range storage.Items {
		if storageItem.Id == req.GetStorageItemId() {
			selectedStorageItem = storageItem
		}
	}

	if selectedStorageItem.Id == "" {
		return nil, status.Error(codes.NotFound, "storage_item not found")
	}

	// Make sure the items are stackable
	if selectedStorageItem.Item.Stackable == false {
		return nil, status.Error(codes.Aborted, "item is not stackable")
	}

	amounts := req.GetAmounts()

	// Use the chunking method to determine the amounts
	if len(amounts) == 0 && req.GetChunks() != 0 {
		// Determine the size of the chunks and add these chunks to the amounts
		chunkSize := math.Floor(float64(selectedStorageItem.Amount) / float64(req.GetChunks()))
		for i := int64(0); i < req.GetChunks(); i++ {
			amounts = append(amounts, int64(chunkSize))
		}

		// Add the remainder to the amounts
		remainderSize := selectedStorageItem.Amount / req.GetChunks()
		if remainderSize > 0 {
			amounts = append(amounts, remainderSize)
		}
	}

	// Use the amount method to determine the amounts
	if len(amounts) == 0 && req.GetAmount() != 0 {
		remainder := selectedStorageItem.Amount

		for remainder > 0 {
			amount := req.GetAmount()
			if amount > remainder {
				amount = remainder
			}

			amounts = append(amounts, amount)
			remainder = remainder - amount
		}
	}

	// Calculate the total amount
	totalAmount := int64(0)
	for _, amount := range amounts {
		totalAmount = totalAmount + amount
	}

	// Amounts dont match
	if totalAmount != selectedStorageItem.Amount {
		return nil, status.Error(codes.Aborted, "The storage_item's amount does not match the given amounts")
	}

	storage, err = s.storageRepository.SplitStack(
		ctx,
		req.GetStorageItemId(),
		amounts,
	)
	if err != nil {
		return nil, status.Error(codes.Aborted, "unable to split stacks")
	}

	return &v1.SplitStackResponse{
		Storage: storage,
	}, nil
}

func (s *economyServiceServer) GiveItem(ctx context.Context, req *v1.GiveItemRequest) (*v1.GiveItemResponse, error) {
	fmt.Println("GiveItem")

	amount := int64(1)

	// Generate a random amount
	if req.GetAmount() != nil {
		amount = random.GenerateRandomInt(
			req.GetAmount().MinAmount,
			req.GetAmount().MaxAmount,
		)
	}

	// Create a remainder so whe know how many items still need to be created
	remainder := amount

	// Get the item
	item, err := s.itemRepository.Get(ctx, req.GetItemId())
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to give item to storage")
	}

	// Increase existing storage_items
	remainder, err = s.GiveToExistingStorageItems(
		ctx,
		req.GetStorageId(),
		req.GetItemId(),
		remainder,
		item,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to give item to storage")
	}

	// Create multiple unstackable items
	if item.Stackable == false && remainder > 0 {
		loops := int(remainder)
		for i := 0; i < loops; i++ {
			_, err := s.storageRepository.GiveItem(
				ctx,
				req.GetStorageId(),
				req.GetItemId(),
				1,
			)
			if err != nil {
				return nil, status.Error(codes.Internal, "unable to give item to storage")
			}

			remainder--
		}
	}

	// For new stacks and full stack for unstackable items
	// and DEFAULT & UNBALANCED_CREATE_NEW_STACKS stack balancing methods
	if item.Stackable == true && remainder > 0 {
		resultAmounts := []int64{}

		if item.StackMaxAmount == 0 {
			resultAmounts = append(resultAmounts, remainder)
		}

		if item.StackMaxAmount > 0 {
			fullStacksToCreate := math.Floor(float64(remainder) / float64(item.StackMaxAmount))

			for i := 0; i < int(fullStacksToCreate); i++ {
				resultAmounts = append(resultAmounts, item.StackMaxAmount)
			}

			resultAmounts = append(resultAmounts, remainder%item.StackMaxAmount)
		}

		for _, resultAmount := range resultAmounts {
			if resultAmount == 0 {
				continue
			}

			_, err := s.storageRepository.GiveItem(
				ctx,
				req.GetStorageId(),
				req.GetItemId(),
				resultAmount,
			)
			if err != nil {
				return nil, status.Error(codes.Internal, "unable to give item to storage")
			}

			remainder -= resultAmount
		}
	}

	if remainder > 0 {
		return nil, status.Error(codes.Internal, "unable to give item to storage")
	}

	return &v1.GiveItemResponse{
		StorageId: req.GetStorageId(),
		Amount:    amount,
	}, nil
}

func (s *economyServiceServer) GetExistingStorageItems(ctx context.Context, storageID string, itemID string) ([]*v1.StorageItem, error) {
	// Get the storage
	storage, err := s.storageRepository.Get(ctx, storageID)
	if err != nil {
		return nil, err
	}

	// Filter out storageItems that are not full
	existingStorageItems := []*v1.StorageItem{}
	for _, storageItem := range storage.Items {
		if storageItem.Item.Id != itemID {
			continue
		}

		if storageItem.Item.StackMaxAmount > 0 && storageItem.Amount >= storageItem.Item.StackMaxAmount {
			continue
		}

		existingStorageItems = append(existingStorageItems, storageItem)
	}

	return existingStorageItems, nil
}

func (s *economyServiceServer) GiveToExistingStorageItems(ctx context.Context, storageID string, itemID string, remainder int64, item *v1.Item) (int64, error) {
	if !item.Stackable {
		return remainder, nil
	}

	// Checks if item is stackable and new items should be added to existing stacks
	if item.StackBalancingMethod != v1.StackBalancingMethod_BALANCED_FILL_EXISTING_STACKS && item.StackBalancingMethod != v1.StackBalancingMethod_UNBALANCED_FILL_EXISTING_STACKS {
		return remainder, nil
	}

	// Get existing storageItems with the same item_id
	existingStorageItems, err := s.GetExistingStorageItems(
		ctx,
		storageID,
		itemID,
	)

	if err != nil {
		return remainder, err
	}

	if len(existingStorageItems) == 0 {
		return remainder, nil
	}

	// An existing stack already exists,
	// It does not have a max_amount so lets increase that one instead
	if item.StackMaxAmount == 0 {
		storageItem := existingStorageItems[0]

		err := s.storageRepository.IncreaseItemAmount(
			ctx,
			storageItem.Id,
			remainder,
		)
		if err != nil {
			return remainder, err
		}

		remainder = 0
	}

	// Because there is a stack_max_amount we should not accidentally overflow it
	// So we'll first try to spread if over the existing stacks
	if item.StackMaxAmount > 0 {
		for _, existingStorageItem := range existingStorageItems {
			// Calculate the remaining space
			existingStorageItemRemainder := item.StackMaxAmount - existingStorageItem.Amount

			// Calculate the amount to increase
			existingStorageItemIncrease := remainder
			if remainder >= existingStorageItemRemainder {
				existingStorageItemIncrease = existingStorageItemRemainder
			}

			err := s.storageRepository.IncreaseItemAmount(
				ctx,
				existingStorageItem.Id,
				existingStorageItemIncrease,
			)
			if err != nil {
				return remainder, err
			}

			remainder -= existingStorageItemIncrease
		}
	}

	return remainder, nil
}
