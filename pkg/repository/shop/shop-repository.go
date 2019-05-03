package shoprepository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"

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
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT
				shop.id as shopId,
				shop.name as shopName,
				shop.created_at as shopCreatedAt,
				shop.updated_at as shopUpdatedAt,
				product.id as productId,
				product.name as productName,
				product.created_at as productCreatedAt,
				product.updated_at as productUpdatedAt,
				item.id as itemId,
				item.name as itemName,
				item.metadata as itemData,
				product_item.id as productItemId,
				product_item.amount as productItemAmount
			FROM shop
      LEFT JOIN shop_product ON (shop_product.shop_id = shop.id)
      LEFT JOIN product ON (product.id = shop_product.product_id)
      LEFT JOIN product_item ON (product_item.product_id = product.id)
      LEFT JOIN item ON (product_item.item_id = item.id)
			WHERE shop.id = $1
		`,
		shopID,
	)

	if err != nil {
		return nil, err
	}

	type row struct {
		ShopID            string
		ShopName          string
		ShopCreatedAt     time.Time
		ShopUpdatedAt     time.Time
		ProductID         sql.NullString
		ProductName       sql.NullString
		ProductCreatedAt  pq.NullTime
		ProductUpdatedAt  pq.NullTime
		ItemID            sql.NullString
		ItemName          sql.NullString
		ItemData          sql.NullString
		ProductItemID     sql.NullString
		ProductItemAmount sql.NullInt64
	}

	productItems := map[string]map[string]*v1.ProductItem{}
	shopProducts := map[string]*v1.Product{}

	var res row
	for rows.Next() {
		err = rows.Scan(
			&res.ShopID,
			&res.ShopName,
			&res.ShopCreatedAt,
			&res.ShopUpdatedAt,
			&res.ProductID,
			&res.ProductName,
			&res.ProductCreatedAt,
			&res.ProductUpdatedAt,
			&res.ItemID,
			&res.ItemName,
			&res.ItemData,
			&res.ProductItemID,
			&res.ProductItemAmount,
		)

		if err != nil {
			return nil, err
		}

		// Create the products
		if res.ProductID.Valid && shopProducts[res.ProductID.String] == nil {
			product := &v1.Product{
				Id:   res.ProductID.String,
				Name: res.ProductName.String,
			}

			product.CreatedAt, _ = ptypes.TimestampProto(res.ProductCreatedAt.Time)
			product.UpdatedAt, _ = ptypes.TimestampProto(res.ProductUpdatedAt.Time)

			shopProducts[res.ProductID.String] = product
		}

		// Extract the Item
		item := &v1.Item{}
		if res.ItemID.Valid && res.ItemName.Valid {
			item.Id = res.ItemID.String
			item.Name = res.ItemName.String
		}

		// Extract the ProductItem
		productItem := &v1.ProductItem{}
		if res.ProductItemID.Valid {
			productItem.Id = res.ProductItemID.String
			productItem.Amount = res.ProductItemAmount.Int64
			productItem.Item = item
		}

		// Add object to the productItems if it is set
		if res.ProductItemID.Valid && res.ItemID.Valid {
			if productItems[res.ProductID.String] == nil {
				productItems[res.ProductID.String] = map[string]*v1.ProductItem{}
			}
			productItems[res.ProductID.String][productItem.Id] = productItem
		}
	}

	// Convert item map into item slice
	products := []*v1.Product{}
	for shopProductKey, shopProductValue := range shopProducts {
		shopProductItems := []*v1.ProductItem{}
		for _, productItemValue := range productItems[shopProductKey] {
			shopProductItems = append(shopProductItems, productItemValue)
		}

		shopProductValue.Items = shopProductItems

		products = append(products, shopProductValue)
	}

	shop := &v1.Shop{
		Id:       res.ShopID,
		Name:     res.ShopName,
		Products: products,
	}

	// Convert created_at to timestamp
	shop.CreatedAt, _ = ptypes.TimestampProto(res.ShopCreatedAt)
	shop.UpdatedAt, _ = ptypes.TimestampProto(res.ShopUpdatedAt)

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
		`DELETE FROM shop_product WHERE product_id = $1 RETURNING shop_id`,
		shopProductID,
	).Scan(&shopID)

	if err != nil {
		return nil, err
	}

	if shopID == "" {
		return nil, fmt.Errorf("unable to retrieve the old shop_id")
	}

	return r.Get(ctx, shopID)
}
