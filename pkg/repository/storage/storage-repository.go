package storagerepository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	repository "github.com/GameComponent/economy-service/pkg/repository"
	jsonpb "github.com/golang/protobuf/jsonpb"
	ptypes "github.com/golang/protobuf/ptypes"
	_struct "github.com/golang/protobuf/ptypes/struct"
	"go.uber.org/zap"
)

// StorageRepository struct
type StorageRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewStorageRepository constructor
func NewStorageRepository(db *sql.DB, logger *zap.Logger) repository.StorageRepository {
	return &StorageRepository{
		db:     db,
		logger: logger,
	}
}

// Create a storage
func (r *StorageRepository) Create(ctx context.Context, playerID string, name string, metadata *_struct.Struct) (*v1.Storage, error) {
	// Parse struct to JSON string
	jsonMetadata := "{}"
	if metadata != nil {
		var err error
		marshaler := jsonpb.Marshaler{}
		jsonMetadata, err = marshaler.MarshalToString(metadata)
		if err != nil {
			return nil, err
		}
	}

	// Add item to the databased return the generated UUID
	lastInsertUUID := ""
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO storage(player_id, name, metadata) VALUES ($1, $2, $3) RETURNING id`,
		playerID,
		name,
		jsonMetadata,
	).Scan(&lastInsertUUID)

	if err != nil {
		return nil, err
	}

	return r.Get(ctx, lastInsertUUID)
}

// Update a storage
func (r *StorageRepository) Update(ctx context.Context, storageID string, name string, metadata *_struct.Struct) (*v1.Storage, error) {
	index := 1
	queries := []string{}
	arguments := []interface{}{}

	// Add name to the query
	if name != "" {
		queries = append(queries, fmt.Sprintf("name = $%v", index))
		arguments = append(arguments, name)
		index++
	}

	// Add metadata to the query
	if metadata != nil {
		// Parse the metadata to a JSON string
		jsonMetadata := "{}"
		var err error
		marshaler := jsonpb.Marshaler{}
		jsonMetadata, err = marshaler.MarshalToString(metadata)
		if err != nil {
			return nil, err
		}

		queries = append(queries, fmt.Sprintf("metadata = $%v", index))
		arguments = append(arguments, jsonMetadata)
		index++
	}

	if index <= 1 {
		return nil, fmt.Errorf("no arguments given")
	}

	// Update the storage
	arguments = append(arguments, storageID)
	query := fmt.Sprintf("UPDATE storage SET %v WHERE id =$%v", strings.Join(queries, ", "), index)
	_, err := r.db.ExecContext(
		ctx,
		query,
		arguments...,
	)
	if err != nil {
		return nil, err
	}

	return r.Get(ctx, storageID)
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
      storage_item.amount as storageItemAmount,
			storage_item.metadata as storageItemData,
      item.id as itemId,
			item.name as itemName,
			item.stackable as itemStackable,
			item.stack_max_amount as itemStackMaxAmount,
			item.stack_balancing_method as itemStackBalancingMethod,
			item.metadata as itemData,
			storage_currency.id as storageCurrencyId,
			storage_currency.amount as storageCurrencyAmount,
      currency.id as currencyId,
      currency.name as currencyName,
      currency.short_name as currencyShortName,
      currency.symbol as currencySymbol
    FROM storage 
    LEFT JOIN storage_item ON (storage.id = storage_item.storage_id)
		LEFT JOIN item ON (storage_item.item_id = item.id)
		LEFT JOIN storage_currency ON (storage.id = storage_currency.storage_id)
    LEFT JOIN currency ON (storage_currency.currency_id = currency.id)
    WHERE storage.id = $1`,
		storageID,
	)

	if err != nil {
		return nil, err
	}

	type row struct {
		StorageID                string
		StorageName              string
		StorageData              string
		PlayerID                 string
		StorageItemID            sql.NullString
		StorageItemAmount        sql.NullInt64
		StorageItemItemData      sql.NullString
		ItemID                   sql.NullString
		ItemName                 sql.NullString
		ItemStackable            sql.NullBool
		ItemStackMaxAmount       sql.NullInt64
		ItemStackBalancingMethod sql.NullInt64
		ItemData                 sql.NullString
		StorageCurrencyID        sql.NullString
		StorageCurrencyAmount    sql.NullInt64
		CurrencyID               sql.NullString
		CurrencyName             sql.NullString
		CurrencyShortName        sql.NullString
		CurrencySymbol           sql.NullString
	}

	storageItems := map[string]*v1.StorageItem{}
	storageCurrencies := map[string]*v1.StorageCurrency{}

	var res row
	for rows.Next() {
		err = rows.Scan(
			&res.StorageID,
			&res.StorageName,
			&res.StorageData,
			&res.PlayerID,
			&res.StorageItemID,
			&res.StorageItemAmount,
			&res.StorageItemItemData,
			&res.ItemID,
			&res.ItemName,
			&res.ItemStackable,
			&res.ItemStackMaxAmount,
			&res.ItemStackBalancingMethod,
			&res.ItemData,
			&res.StorageCurrencyID,
			&res.StorageCurrencyAmount,
			&res.CurrencyID,
			&res.CurrencyName,
			&res.CurrencyShortName,
			&res.CurrencySymbol,
		)

		if err != nil {
			return nil, err
		}

		// Extract the Item
		item := &v1.Item{}
		if res.ItemID.Valid && res.ItemName.Valid {
			item.Id = res.ItemID.String
			item.Name = res.ItemName.String
			item.Stackable = res.ItemStackable.Bool
			item.StackMaxAmount = res.ItemStackMaxAmount.Int64
			item.StackBalancingMethod = v1.StackBalancingMethod(res.ItemStackBalancingMethod.Int64)
		}

		// Extract the StorageItem
		storageItem := &v1.StorageItem{}
		if res.StorageItemID.Valid {
			storageItem.Id = res.StorageItemID.String
			storageItem.Item = item

			// Only show the amount if the item is stackable
			if res.ItemStackable.Bool {
				storageItem.Amount = res.StorageItemAmount.Int64
			}
		}

		// Extract the Currency
		currency := &v1.Currency{}
		if res.CurrencyID.Valid && res.CurrencyName.Valid {
			currency.Id = res.CurrencyID.String
			currency.Name = res.CurrencyName.String
			currency.ShortName = res.CurrencyShortName.String
			currency.Symbol = res.CurrencySymbol.String
		}

		// Extract the StorageCurrency
		storageCurrency := &v1.StorageCurrency{}
		if res.StorageCurrencyID.Valid {
			storageCurrency.Id = res.StorageCurrencyID.String
			storageCurrency.Amount = res.StorageCurrencyAmount.Int64
			storageCurrency.Currency = currency
		}

		// Add object to the storageItem if it is set
		if storageItem.Id != "" {
			storageItems[storageItem.Id] = storageItem
		}

		// Add object to the storageCurrency if it is set
		if storageCurrency.Id != "" {
			storageCurrencies[storageCurrency.Id] = storageCurrency
		}
	}

	// Convert item map into item slice
	items := []*v1.StorageItem{}
	for _, value := range storageItems {
		items = append(items, value)
	}

	// Convert currency map into currency slice
	currencies := []*v1.StorageCurrency{}
	for _, value := range storageCurrencies {
		currencies = append(currencies, value)
	}

	// Convert metadata json to a proto Struct
	stringReader := strings.NewReader(res.StorageData)
	metadataStruct := _struct.Struct{}
	unmarshaler := jsonpb.Unmarshaler{}
	err = unmarshaler.Unmarshal(stringReader, &metadataStruct)
	if err != nil {
		return nil, err
	}

	storage := &v1.Storage{
		Id:         res.StorageID,
		PlayerId:   res.PlayerID,
		Name:       res.StorageName,
		Items:      items,
		Currencies: currencies,
		Metadata:   &metadataStruct,
	}

	return storage, nil
}

// GiveItem to a storage
func (r *StorageRepository) GiveItem(ctx context.Context, storageID string, itemID string, amount int64) (*string, error) {
	// Add item to the databased return the generated UUID
	lastInsertUUID := ""
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO storage_item(item_id, storage_id, amount) VALUES ($1, $2, $3) RETURNING id`,
		itemID,
		storageID,
		amount,
	).Scan(&lastInsertUUID)

	if err != nil {
		return nil, err
	}

	return &lastInsertUUID, nil
}

// IncreaseItemAmount to a storage
func (r *StorageRepository) IncreaseItemAmount(ctx context.Context, storageItemID string, amount int64) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE storage_item SET amount = amount + $1 WHERE id = $2`,
		amount,
		storageItemID,
	)

	if err != nil {
		return err
	}

	return nil
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
		Id:     storageCurrencyUUID,
		Amount: storageCurrencyAmount,
	}

	return storageCurrency, nil
}

// List all storages
func (r *StorageRepository) List(ctx context.Context, limit int32, offset int32) ([]*v1.Storage, int32, error) {
	// Query items from the database
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT 
				id,
				name,
				player_id,
				created_at,
				updated_at,
				(SELECT COUNT(DISTINCT id) FROM storage) AS total_size
			FROM storage
			LIMIT $1
			OFFSET $2
		`,
		limit,
		offset,
	)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Unwrap rows into items
	storages := []*v1.Storage{}
	totalSize := int32(0)

	for rows.Next() {
		storage := v1.Storage{}
		createdAt := time.Time{}
		updatedAt := time.Time{}

		err := rows.Scan(
			&storage.Id,
			&storage.Name,
			&storage.PlayerId,
			&createdAt,
			&updatedAt,
			&totalSize,
		)
		if err != nil {
			return nil, 0, err
		}

		// Convert created_at to timestamp
		storage.CreatedAt, _ = ptypes.TimestampProto(createdAt)
		storage.UpdatedAt, _ = ptypes.TimestampProto(updatedAt)

		storages = append(storages, &storage)
	}

	return storages, totalSize, nil
}
