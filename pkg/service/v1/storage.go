package v1

import (
	"context"
	"fmt"
	"log"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	"github.com/google/uuid"
)

func (s *economyServiceServer) GiveItem(ctx context.Context, req *v1.GiveItemRequest) (*v1.GiveItemResponse, error) {
	fmt.Println("GiveItem")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// Add item to the databased return the generated UUID
	lastInsertUUID := ""
	err := s.db.QueryRowContext(
		ctx,
		`INSERT INTO storage_item(item_id, storage_id) VALUES ($1, $2) RETURNING id`,
		req.GetItemId(),
		req.GetStorageId(),
	).Scan(&lastInsertUUID)

	if err != nil {
		return nil, err
	}

	item := &v1.Item{}

	err = s.db.QueryRowContext(
		ctx,
		`SELECT id, name FROM item WHERE id = $1`,
		req.GetItemId(),
	).Scan(&item.Id, &item.Name)

	if err != nil {
		return nil, err
	}

	storageItem := &v1.StorageItem{
		Id:   lastInsertUUID,
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

	// Add item to the databased return the generated UUID
	lastInsertUUID := ""
	err := s.db.QueryRowContext(
		ctx,
		`INSERT INTO storage(player_id, name) VALUES ($1, $2) RETURNING id`,
		req.GetPlayerId(),
		req.GetName(),
	).Scan(&lastInsertUUID)

	if err != nil {
		return nil, err
	}

	storage := &v1.Storage{
		Id:       lastInsertUUID,
		PlayerId: req.GetPlayerId(),
		Name:     req.GetName(),
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

	rows, err := s.db.QueryContext(
		ctx,
		`SELECT
      storage.id as storageId,
      storage.name as storageName,
      storage.metadata as storageData,
      storage.player_id as playerId,
      storage_item.id as storageItemId,
      storage_item.metadata as storageItemData,
      item.id as itemId,
      item.name as itemName,
      item.metadata as itemData
    FROM storage 
    INNER JOIN storage_item on (storage.id = storage_item.storage_id)
    INNER JOIN item on (storage_item.item_id = item.id)
    WHERE storage.id = $1`,
		req.GetStorageId(),
	)

	if err != nil {
		return nil, err
	}

	storageItems := []*v1.StorageItem{}

	type row struct {
		StorageID           string
		StorageName         string
		StorageData         string
		PlayerID            string
		StorageItemID       string
		StorageItemItemData string
		ItemID              string
		ItemName            string
		ItemData            string
	}

	var r row
	for rows.Next() {
		err = rows.Scan(
			&r.StorageID,
			&r.StorageName,
			&r.StorageData,
			&r.PlayerID,
			&r.StorageItemID,
			&r.StorageItemItemData,
			&r.ItemID,
			&r.ItemName,
			&r.ItemData,
		)
		if err != nil {
			log.Fatal(err)
		}

		item := &v1.Item{
			Id:   r.ItemID,
			Name: r.ItemName,
		}

		storageItem := &v1.StorageItem{
			Id:   r.StorageItemID,
			Item: item,
		}

		storageItems = append(storageItems, storageItem)
	}

	storage := &v1.Storage{
		Id:       r.StorageID,
		PlayerId: r.PlayerID,
		Name:     r.StorageName,
		Items:    storageItems,
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

	// Add item to the databased return the generated UUID
	storageCurrencyUuid := ""
	storageCurrencyAmount := int64(0)

	err := s.db.QueryRowContext(
		ctx, `
      INSERT INTO storage_currency(currency_id, storage_id, amount)
      VALUES($1, $2, $3)
      ON CONFLICT(currency_id,storage_id) DO UPDATE
      SET amount = storage_currency.amount + EXCLUDED.amount
      RETURNING id, amount
    `,
		req.GetCurrencyId(),
		req.GetStorageId(),
		req.GetAmount(),
	).Scan(&storageCurrencyUuid, &storageCurrencyAmount)

	if err != nil {
		return nil, err
	}

	storageCurrency := &v1.StorageCurrency{
		Id:         storageCurrencyUuid,
		CurrencyId: req.GetCurrencyId(),
		Amount:     storageCurrencyAmount,
	}

	return &v1.GiveCurrencyResponse{
		Api:      apiVersion,
		Currency: storageCurrency,
	}, nil
}
