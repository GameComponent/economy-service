package itemrepository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	repository "github.com/GameComponent/economy-service/pkg/repository"
	jsonpb "github.com/golang/protobuf/jsonpb"
	ptypes "github.com/golang/protobuf/ptypes"
	_struct "github.com/golang/protobuf/ptypes/struct"
)

// ItemRepository struct
type ItemRepository struct {
	db *sql.DB
}

// NewItemRepository constructor
func NewItemRepository(db *sql.DB) repository.ItemRepository {
	return &ItemRepository{
		db: db,
	}
}

// Create a new item
func (r *ItemRepository) Create(ctx context.Context, name string, stackable bool, stackMaxAmount int64, stackBalancingMethod int64, metadata *_struct.Struct) (*v1.Item, error) {
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

	lastInsertUUID := ""
	err := r.db.QueryRowContext(
		ctx,
		`
			INSERT INTO item(
				name,
				stackable,
				stack_max_amount,
				stack_balancing_method,
				metadata
			)
			VALUES (
				$1,
				$2,
				$3,
				$4,
				$5
			)
			RETURNING id
		`,
		name,
		stackable,
		stackMaxAmount,
		stackBalancingMethod,
		jsonMetadata,
	).Scan(&lastInsertUUID)

	if err != nil {
		return nil, err
	}

	return r.Get(ctx, lastInsertUUID)
}

// Update an item
func (r *ItemRepository) Update(ctx context.Context, id string, name string, metadata *_struct.Struct) (*v1.Item, error) {
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

	_, err := r.db.ExecContext(
		ctx,
		`UPDATE item SET name = $1, metadata = $2 WHERE id = $3`,
		name,
		jsonMetadata,
		id,
	)

	if err != nil {
		return nil, err
	}

	return r.Get(ctx, id)
}

// List all items
func (r *ItemRepository) List(ctx context.Context, limit int32, offset int32) ([]*v1.Item, int32, error) {
	// Query items from the database
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT 
				id,
				name,
				created_at,
				updated_at,
				(SELECT COUNT(*) FROM item) AS total_size
			FROM item
			ORDER BY created_at DESC 
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
	items := []*v1.Item{}
	totalSize := int32(0)

	for rows.Next() {
		item := v1.Item{}
		createdAt := time.Time{}
		updatedAt := time.Time{}

		err := rows.Scan(
			&item.Id,
			&item.Name,
			&createdAt,
			&updatedAt,
			&totalSize,
		)
		if err != nil {
			return nil, 0, err
		}

		// Convert created_at to timestamp
		item.CreatedAt, _ = ptypes.TimestampProto(createdAt)
		item.UpdatedAt, _ = ptypes.TimestampProto(updatedAt)

		items = append(items, &item)
	}

	return items, totalSize, nil
}

// Get an item
func (r *ItemRepository) Get(ctx context.Context, itemID string) (*v1.Item, error) {
	item := &v1.Item{}
	createdAt := time.Time{}
	updatedAt := time.Time{}
	jsonMetadata := ""

	err := r.db.QueryRowContext(
		ctx,
		`
			SELECT
				id,
				name,
				stackable,
				stack_max_amount,
				stack_balancing_method,
				created_at,
				updated_at,
				metadata
			FROM item
			WHERE id = $1
		`,
		itemID,
	).Scan(
		&item.Id,
		&item.Name,
		&item.Stackable,
		&item.StackMaxAmount,
		&item.StackBalancingMethod,
		&createdAt,
		&updatedAt,
		&jsonMetadata,
	)
	if err != nil {
		return nil, err
	}

	// Convert metadata json to a proto Struct
	stringReader := strings.NewReader(jsonMetadata)
	metadataStruct := _struct.Struct{}
	unmarshaler := jsonpb.Unmarshaler{}
	err = unmarshaler.Unmarshal(stringReader, &metadataStruct)
	if err != nil {
		return nil, err
	}
	item.Metadata = &metadataStruct

	// Convert created_at to timestamp
	item.CreatedAt, _ = ptypes.TimestampProto(createdAt)
	item.UpdatedAt, _ = ptypes.TimestampProto(updatedAt)

	return item, nil
}

// Search item
func (r *ItemRepository) Search(ctx context.Context, query string, limit int32, offset int32) ([]*v1.Item, int32, error) {
	// Query items from the database
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT
				id,
				name,
				created_at,
				updated_at,
				(SELECT COUNT(id) FROM item WHERE name ~* $1) AS total_size
			FROM item
			WHERE name ~* $1
			LIMIT $2
			OFFSET $3
		`,
		query,
		limit,
		offset,
	)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Unwrap rows into items
	items := []*v1.Item{}
	totalSize := int32(0)

	for rows.Next() {
		item := v1.Item{}
		createdAt := time.Time{}
		updatedAt := time.Time{}

		err := rows.Scan(
			&item.Id,
			&item.Name,
			&createdAt,
			&updatedAt,
			&totalSize,
		)
		if err != nil {
			return nil, 0, err
		}

		// Convert created_at to timestamp
		item.CreatedAt, _ = ptypes.TimestampProto(createdAt)
		item.UpdatedAt, _ = ptypes.TimestampProto(updatedAt)

		items = append(items, &item)
	}

	return items, totalSize, nil
}
