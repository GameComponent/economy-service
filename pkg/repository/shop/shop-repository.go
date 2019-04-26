package shoprepository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	"github.com/golang/protobuf/ptypes"
)

// ShopRepository struct
type ShopRepository struct {
	db *sql.DB
}

// NewShopRepository constructor
func NewShopRepository(db *sql.DB) *ShopRepository {
	return &ShopRepository{
		db: db,
	}
}

// Get a shop
func (r *ShopRepository) Get(ctx context.Context, shopID string) (*v1.Shop, error) {
	shop := &v1.Shop{}
	createdAt := time.Time{}
	updatedAt := time.Time{}

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, created_at, updated_at FROM shop WHERE id = $1`,
		shopID,
	).Scan(
		&shop.Id,
		&shop.Name,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Convert created_at to timestamp
	shop.CreatedAt, _ = ptypes.TimestampProto(createdAt)
	shop.UpdatedAt, _ = ptypes.TimestampProto(updatedAt)

	return shop, nil
}

// Create a new shop
func (r *ShopRepository) Create(ctx context.Context, name string) (*v1.Shop, error) {
	databaseID := ""
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO shop(name) VALUES ($1) RETURNING id`,
		name,
	).Scan(&databaseID)

	if err != nil {
		return nil, err
	}

	return &v1.Shop{
		Id:   databaseID,
		Name: name,
	}, nil
}

// List all shops
func (r *ShopRepository) List(
	ctx context.Context,
	limit int32,
	offset int32,
) (
	[]*v1.Shop,
	int32,
	error,
) {
	// Query shops from the database
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT
				id,
				name,
				(SELECT COUNT(id) FROM shop) AS total_size
			FROM shop
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

	// Unwrap rows into shops
	shops := []*v1.Shop{}
	totalSize := int32(1)

	for rows.Next() {
		shop := v1.Shop{}

		err := rows.Scan(
			&shop.Id,
			&shop.Name,
			&totalSize,
		)
		if err != nil {
			return nil, 0, err
		}

		shops = append(shops, &shop)
	}

	return shops, totalSize, nil
}

// AttachProduct to a shop
func (r *ShopRepository) AttachProduct(ctx context.Context, shopID string, productID string) (*v1.Shop, error) {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO shop_product(shop_id, product_id) VALUES ($1, $2)`,
		shopID,
		productID,
	)

	if err != nil {
		return nil, err
	}

	return r.Get(ctx, shopID)
}

// DetachProduct from a shop
func (r *ShopRepository) DetachProduct(ctx context.Context, shopProductID string) (*v1.Shop, error) {
	shopID := ""
	err := r.db.QueryRowContext(
		ctx,
		`DELETE FROM shop_product WHERE id = $1 RETURNING shop_id`,
		shopProductID,
	).Scan(&shopID)

	if err != nil {
		return nil, err
	}

	if shopID == "" {
		return nil, fmt.Errorf("unable to retrieve the old product_id")
	}

	return r.Get(ctx, shopID)
}
