package productrepository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
)

// ProductRepository struct
type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository constructor
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

// Create a new product
func (r *ProductRepository) Create(ctx context.Context, name string) (*v1.Product, error) {
	lastInsertUUID := ""
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO product(name) VALUES ($1) RETURNING id`,
		name,
	).Scan(&lastInsertUUID)

	if err != nil {
		return nil, err
	}

	return &v1.Product{
		Id:   lastInsertUUID,
		Name: name,
	}, nil
}

// Update an product
func (r *ProductRepository) Update(ctx context.Context, id string, name string) (*v1.Product, error) {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE product SET name = $1 WHERE id = $3`,
		name,
		id,
	)

	if err != nil {
		return nil, err
	}

	product := &v1.Product{
		Id:   id,
		Name: name,
	}

	return product, nil
}

// List all products
func (r *ProductRepository) List(
	ctx context.Context,
	limit int32,
	offset int32,
) (
	[]*v1.Product,
	int32,
	error,
) {
	// Query products from the database
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT 
				id,
				name,
				created_at,
				updated_at,
				(SELECT COUNT(*) FROM product) AS total_size
			FROM product
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

	// Unwrap rows into products
	products := []*v1.Product{}
	totalSize := int32(0)

	for rows.Next() {
		product := v1.Product{}
		createdAt := time.Time{}
		updatedAt := time.Time{}

		err := rows.Scan(
			&product.Id,
			&product.Name,
			&createdAt,
			&updatedAt,
			&totalSize,
		)
		if err != nil {
			return nil, 0, err
		}

		// Convert created_at to timestamp
		product.CreatedAt, _ = ptypes.TimestampProto(createdAt)
		product.UpdatedAt, _ = ptypes.TimestampProto(updatedAt)

		products = append(products, &product)
	}

	return products, totalSize, nil
}

// Get an product
func (r *ProductRepository) Get(ctx context.Context, productID string) (*v1.Product, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT
				product.id as productId,
				product.name as productName,
				product.created_at as productCreatedAt,
				product.updated_at as productUpdatedAt,
				item.id as itemId,
				item.name as itemName,
				item.metadata as itemData,
				product_item.id as productItemId,
				product_item.amount as productItemAmount
			FROM product
			LEFT JOIN product_item ON (product_item.product_id = product.id)
			LEFT JOIN item ON (item.id = product_item.item_id)
			WHERE product.id = $1
		`,
		productID,
	)

	if err != nil {
		return nil, err
	}

	type row struct {
		ProductID         string
		ProductName       string
		ProductCreatedAt  time.Time
		ProductUpdatedAt  time.Time
		ItemID            sql.NullString
		ItemName          sql.NullString
		ItemData          sql.NullString
		ProductItemID     sql.NullString
		ProductItemAmount sql.NullInt64
	}

	productItems := map[string]*v1.ProductItem{}

	var res row
	for rows.Next() {
		err = rows.Scan(
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
		if productItem.Item != nil {
			productItems[productItem.Id] = productItem
		}
	}

	// Convert item map into item slice
	items := []*v1.ProductItem{}
	for _, value := range productItems {
		items = append(items, value)
	}

	product := &v1.Product{
		Id:    res.ProductID,
		Name:  res.ProductName,
		Items: items,
	}

	// Convert created_at to timestamp
	product.CreatedAt, _ = ptypes.TimestampProto(res.ProductCreatedAt)
	product.UpdatedAt, _ = ptypes.TimestampProto(res.ProductUpdatedAt)

	return product, nil
}

// Search product
func (r *ProductRepository) Search(
	ctx context.Context,
	query string,
	limit int32,
	offset int32,
) (
	[]*v1.Product,
	int32,
	error,
) {
	// Query products from the database
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT
				id,
				name,
				created_at,
				updated_at,
				(SELECT COUNT(id) FROM product WHERE name ~* $1) AS total_size
			FROM product
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

	// Unwrap rows into products
	products := []*v1.Product{}
	totalSize := int32(0)

	for rows.Next() {
		product := v1.Product{}
		createdAt := time.Time{}
		updatedAt := time.Time{}

		err := rows.Scan(
			&product.Id,
			&product.Name,
			&createdAt,
			&updatedAt,
			&totalSize,
		)
		if err != nil {
			return nil, 0, err
		}

		// Convert created_at to timestamp
		product.CreatedAt, _ = ptypes.TimestampProto(createdAt)
		product.UpdatedAt, _ = ptypes.TimestampProto(updatedAt)

		products = append(products, &product)
	}

	return products, totalSize, nil
}

// AttachItem to a product
func (r *ProductRepository) AttachItem(ctx context.Context, productID string, itemID string, amount int64) (*v1.Product, error) {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO product_item(product_id, item_id, amount) VALUES ($1, $2, $3)`,
		productID,
		itemID,
		amount,
	)

	if err != nil {
		return nil, err
	}

	return r.Get(ctx, productID)
}

// DetachItem from a product
func (r *ProductRepository) DetachItem(ctx context.Context, productItemID string) (*v1.Product, error) {
	productID := ""
	err := r.db.QueryRowContext(
		ctx,
		`DELETE FROM product_item WHERE id = $1 RETURNING product_id`,
		productItemID,
	).Scan(&productID)

	if err != nil {
		return nil, err
	}

	if productID == "" {
		return nil, fmt.Errorf("unable to retrieve the old product_id")
	}

	return r.Get(ctx, productID)
}

// ListPrice for the product
func (r *ProductRepository) ListPrice(
	ctx context.Context,
	productID string,
) (
	[]*v1.Price,
	error,
) {
	// Query products from the database
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT 
				id,
				created_at,
				updated_at
			FROM price
			ORDER BY created_at DESC 
		`,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Unwrap rows into prices
	prices := []*v1.Price{}

	for rows.Next() {
		price := v1.Price{}
		createdAt := time.Time{}
		updatedAt := time.Time{}

		err := rows.Scan(
			&price.Id,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Convert created_at to timestamp
		price.CreatedAt, _ = ptypes.TimestampProto(createdAt)
		price.UpdatedAt, _ = ptypes.TimestampProto(updatedAt)

		prices = append(prices, &price)
	}

	return prices, nil
}
