package storagerepository

import (
	"context"
	"database/sql"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
)

// StorageRepository struct
type StorageRepository struct {
	db *sql.DB
}

// NewStorageRepository constructor
func NewStorageRepository(db *sql.DB) *StorageRepository {
	return &StorageRepository{
		db: db,
	}
}

// Create a storage
func (r *StorageRepository) Create(ctx context.Context, playerID string, name string) (*v1.Storage, error) {
	// Add item to the databased return the generated UUID
	lastInsertUUID := ""
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO storage(player_id, name) VALUES ($1, $2) RETURNING id`,
		playerID,
		name,
	).Scan(&lastInsertUUID)

	if err != nil {
		return nil, err
	}

	storage := &v1.Storage{
		Id:       lastInsertUUID,
		PlayerId: playerID,
		Name:     name,
	}

	return storage, nil
}

// Get a storage
func (r *StorageRepository) Get(ctx context.Context, storageID string) (*v1.Storage, error) {
	rows, err := r.db.QueryContext(
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
		storageID,
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

	var res row
	for rows.Next() {
		err = rows.Scan(
			&res.StorageID,
			&res.StorageName,
			&res.StorageData,
			&res.PlayerID,
			&res.StorageItemID,
			&res.StorageItemItemData,
			&res.ItemID,
			&res.ItemName,
			&res.ItemData,
		)

		if err != nil {
			return nil, err
		}

		item := &v1.Item{
			Id:   res.ItemID,
			Name: res.ItemName,
		}

		storageItem := &v1.StorageItem{
			Id:   res.StorageItemID,
			Item: item,
		}

		storageItems = append(storageItems, storageItem)
	}

	storage := &v1.Storage{
		Id:       res.StorageID,
		PlayerId: res.PlayerID,
		Name:     res.StorageName,
		Items:    storageItems,
	}

	return storage, nil
}

// GiveItem to a storage
func (r *StorageRepository) GiveItem(ctx context.Context, storageID string, itemID string) (*string, error) {
	// Add item to the databased return the generated UUID
	lastInsertUUID := ""
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO storage_item(item_id, storage_id) VALUES ($1, $2) RETURNING id`,
		itemID,
		storageID,
	).Scan(&lastInsertUUID)

	if err != nil {
		return nil, err
	}

	return &lastInsertUUID, nil
}

// GiveCurrency to a storage
func (r *StorageRepository) GiveCurrency(ctx context.Context, storageID string, currencyID string, amount int64) (*v1.StorageCurrency, error) {
	// Add item to the databased return the generated UUID
	storageCurrencyUUID := ""
	storageCurrencyAmount := int64(0)

	err := r.db.QueryRowContext(
		ctx,
		`
      INSERT INTO storage_currency(currency_id, storage_id, amount)
      VALUES($1, $2, $3)
      ON CONFLICT(currency_id,storage_id) DO UPDATE
      SET amount = storage_currency.amount + EXCLUDED.amount
      RETURNING id, amount
    `,
		currencyID,
		storageID,
		amount,
	).Scan(&storageCurrencyUUID, &storageCurrencyAmount)

	if err != nil {
		return nil, err
	}

	storageCurrency := &v1.StorageCurrency{
		Id:         storageCurrencyUUID,
		CurrencyId: currencyID,
		Amount:     storageCurrencyAmount,
	}

	return storageCurrency, nil
}
