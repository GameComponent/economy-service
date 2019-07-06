package itemrepository

import (
	"context"
	"database/sql"
	"time"

	"github.com/golang/protobuf/ptypes"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	repository "github.com/GameComponent/economy-service/pkg/repository"
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
func (r *ItemRepository) Create(ctx context.Context, name string, stackable bool, stackMaxAmount int64, stackBalancingMethod int64) (*v1.Item, error) {
	lastInsertUUID := ""
	err := r.db.QueryRowContext(
		ctx,
		`
			INSERT INTO item(
				name,
				stackable,
				stack_max_amount,
				stack_balancing_method
			)
			VALUES (
				$1,
				$2,
				$3,
				$4
			)
			RETURNING id
		`,
		name,
		stackable,
		stackMaxAmount,
		stackBalancingMethod,
	).Scan(&lastInsertUUID)

	if err != nil {
		return nil, err
	}

	return &v1.Item{
		Id:   lastInsertUUID,
		Name: name,
	}, nil
}

// Update an item
func (r *ItemRepository) Update(ctx context.Context, id string, name string, metadata string) (*v1.Item, error) {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE item SET name = $1, metadata = $2 WHERE id = $3`,
		name,
		metadata,
		id,
	)

	if err != nil {
		return nil, err
	}

	item := &v1.Item{
		Id:   id,
		Name: name,
	}

	return item, nil
}

// List all items
func (r *ItemRepository) List(ctx context.Context, limit int32, offset int32) ([]*v1.Item, int32, error ) {
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
				updated_at
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
	)
	if err != nil {
		return nil, err
	}

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
