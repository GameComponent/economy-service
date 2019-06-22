package productrepository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
)

// ProductRepository struct
type ProductRepository struct {
	db *sql.DB
}

// NullTime is a nullable time.Time
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
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
				product.id AS productId,
				product.name AS productName,
				product.created_at AS productCreatedAt,
				product.updated_at AS productUpdatedAt,
				product_item.id AS productItemId,
				product_item.amount AS productItemAmount,
				product_currency.id AS productCurrencyId,
				product_currency.amount AS productCurrencyAmount,
				item.id AS itemId,
				item.name AS itemName,
				item.stackable AS itemtackable,
				item.stack_max_amount AS itemStackMaxAmount,
				item.stack_balancing_method AS itemStackBalancingMethod,
				item.created_at AS itemCreatedAt,
				item.updated_at AS itemUpdatedAt,
				currency.id AS currencyId,
				currency.name AS currencyName,
				currency.short_name AS currencyShortName,
				currency.symbol AS currencySymbol,
				price.id AS priceId,
				price_currency.id AS priceCurrecyId,
				price_currency.amount AS priceCurrencyAmount,
				price_item.id AS priceItemId,
				price_item.amount AS priceItemAmount,
				price_currency_currency.id AS priceCurrencyCurrencyId,
				price_currency_currency.name AS priceCurrencyCurrencyName,
				price_currency_currency.short_name AS priceCurrencyCurrencyShortName,
				price_currency_currency.symbol AS priceCurrencyCurrencySymbol,
				price_item_item.id AS priceItemItemId,
				price_item_item.name AS priceItemItemName,
				price_item_item.stackable AS priceItemItemtackable,
				price_item_item.stack_max_amount AS priceItemItemStackMaxAmount,
				price_item_item.stack_balancing_method AS priceItemItemStackBalancingMethod,
				price_item_item.created_at AS priceItemItemCreatedAt,
				price_item_item.updated_at AS priceItemItemUpdatedAt
			FROM product
			LEFT JOIN product_item ON (product_item.product_id = product.id)
			LEFT JOIN item ON (item.id = product_item.item_id)
			LEFT JOIN product_currency ON (product_currency.product_id = product.id)
			LEFT JOIN currency ON (currency.id = product_currency.currency_id)
			LEFT JOIN price ON (price.product_id = product.id)
			LEFT JOIN price_currency ON (price_currency.price_id = price.id)
			LEFT JOIN price_item ON (price_item.price_id = price.id)
			LEFT JOIN currency price_currency_currency ON (price_currency_currency.id = price_currency.currency_id)
			LEFT JOIN item price_item_item ON (price_item_item.id = price_item.item_id)
			WHERE product.id = $1
		`,
		productID,
	)

	if err != nil {
		return nil, err
	}

	type row struct {
		ProductID                         string
		ProductName                       string
		ProductCreatedAt                  time.Time
		ProductUpdatedAt                  time.Time
		ProductItemID                     sql.NullString
		ProductItemAmount                 sql.NullInt64
		ProductCurrencyID                 sql.NullString
		ProductCurrencyAmount             sql.NullInt64
		ItemID                            sql.NullString
		ItemName                          sql.NullString
		ItemStackable                     sql.NullBool
		ItemStackMaxAmount                sql.NullInt64
		ItemStackBalancingMethod          sql.NullInt64
		ItemCreatedAt                     NullTime
		ItemUpdatedAt                     NullTime
		CurrencyID                        sql.NullString
		CurrencyName                      sql.NullString
		CurrencyShortName                 sql.NullString
		CurrencySymbol                    sql.NullString
		PriceID                           sql.NullString
		PriceCurrecyID                    sql.NullString
		PriceCurrencyAmount               sql.NullInt64
		PriceItemID                       sql.NullString
		PriceItemAmount                   sql.NullInt64
		PriceCurrencyCurrencyID           sql.NullString
		PriceCurrencyCurrencyName         sql.NullString
		PriceCurrencyCurrencyShortName    sql.NullString
		PriceCurrencyCurrencySymbol       sql.NullString
		PriceItemItemID                   sql.NullString
		PriceItemItemName                 sql.NullString
		PriceItemItemStackable            sql.NullBool
		PriceItemItemStackMaxAmount       sql.NullInt64
		PriceItemItemStackBalancingMethod sql.NullInt64
		PriceItemItemCreatedAt            NullTime
		PriceItemItemUpdatedAt            NullTime
	}

	productItems := map[string]*v1.ProductItem{}
	productCurrencies := map[string]*v1.ProductCurrency{}
	productPrices := map[string]*v1.Price{}
	productPriceItems := map[string]map[string]*v1.PriceItem{}
	productPriceCurrencies := map[string]map[string]*v1.PriceCurrency{}

	var res row
	for rows.Next() {
		err = rows.Scan(
			&res.ProductID,
			&res.ProductName,
			&res.ProductCreatedAt,
			&res.ProductUpdatedAt,
			&res.ProductItemID,
			&res.ProductItemAmount,
			&res.ProductCurrencyID,
			&res.ProductCurrencyAmount,
			&res.ItemID,
			&res.ItemName,
			&res.ItemStackable,
			&res.ItemStackMaxAmount,
			&res.ItemStackBalancingMethod,
			&res.ItemCreatedAt,
			&res.ItemUpdatedAt,
			&res.CurrencyID,
			&res.CurrencyName,
			&res.CurrencyShortName,
			&res.CurrencySymbol,
			&res.PriceID,
			&res.PriceCurrecyID,
			&res.PriceCurrencyAmount,
			&res.PriceItemID,
			&res.PriceItemAmount,
			&res.PriceCurrencyCurrencyID,
			&res.PriceCurrencyCurrencyName,
			&res.PriceCurrencyCurrencyShortName,
			&res.PriceCurrencyCurrencySymbol,
			&res.PriceItemItemID,
			&res.PriceItemItemName,
			&res.PriceItemItemStackable,
			&res.PriceItemItemStackMaxAmount,
			&res.PriceItemItemStackBalancingMethod,
			&res.PriceItemItemCreatedAt,
			&res.PriceItemItemUpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		// Extract the Item
		item := &v1.Item{}
		if res.ItemID.Valid {
			item.Id = res.ItemID.String
			item.Name = res.ItemName.String
			item.Stackable = res.PriceItemItemStackable.Bool
			item.StackMaxAmount = res.PriceItemItemStackMaxAmount.Int64
			item.StackBalancingMethod = v1.StackBalancingMethod(res.PriceItemItemStackBalancingMethod.Int64)
			item.CreatedAt, _ = ptypes.TimestampProto(res.ItemCreatedAt.Time)
			item.UpdatedAt, _ = ptypes.TimestampProto(res.ItemUpdatedAt.Time)
		}

		// Extract the ProductItem
		productItem := &v1.ProductItem{}
		if res.ProductItemID.Valid {
			productItem.Id = res.ProductItemID.String
			productItem.Amount = res.ProductItemAmount.Int64
			productItem.Item = item
		}

		// Extract the Item
		currency := &v1.Currency{}
		if res.CurrencyID.Valid {
			currency.Id = res.CurrencyID.String
			currency.Name = res.CurrencyName.String
			currency.ShortName = res.CurrencyShortName.String
			currency.Symbol = res.CurrencySymbol.String
		}

		// Extract the ProductCurrency
		productCurrency := &v1.ProductCurrency{}
		if res.ProductCurrencyID.Valid {
			productCurrency.Id = res.ProductCurrencyID.String
			productCurrency.Amount = res.ProductCurrencyAmount.Int64
			productCurrency.Currency = currency
		}

		// Extract Price
		price := &v1.Price{}
		if res.PriceID.Valid {
			price.Id = res.PriceID.String
		}

		// Extract PriceCurrency
		priceCurrency := &v1.PriceCurrency{}
		if res.PriceCurrecyID.Valid {
			priceCurrency.Id = res.PriceCurrecyID.String
			priceCurrency.Amount = res.PriceCurrencyAmount.Int64
		}

		// Extract Currency
		priceCurrencyCurrency := &v1.Currency{}
		if res.PriceCurrencyCurrencyID.Valid {
			priceCurrencyCurrency.Id = res.PriceCurrencyCurrencyID.String
			priceCurrencyCurrency.Name = res.PriceCurrencyCurrencyName.String
			priceCurrencyCurrency.ShortName = res.PriceCurrencyCurrencyShortName.String
			priceCurrencyCurrency.Symbol = res.PriceCurrencyCurrencySymbol.String
		}

		// Extract PriceItem
		priceItem := &v1.PriceItem{}
		if res.PriceItemID.Valid {
			priceItem.Id = res.PriceItemID.String
			priceItem.Amount = res.PriceItemAmount.Int64
		}

		// Extract Item
		priceItemItem := &v1.Item{}
		if res.PriceItemItemID.Valid {
			priceItemItem.Id = res.PriceItemItemID.String
			priceItemItem.Name = res.PriceItemItemName.String
			priceItemItem.Stackable = res.PriceItemItemStackable.Bool
			priceItemItem.StackMaxAmount = res.PriceItemItemStackMaxAmount.Int64
			priceItemItem.StackBalancingMethod = v1.StackBalancingMethod(res.PriceItemItemStackBalancingMethod.Int64)
			priceItemItem.CreatedAt, _ = ptypes.TimestampProto(res.PriceItemItemCreatedAt.Time)
			priceItemItem.UpdatedAt, _ = ptypes.TimestampProto(res.PriceItemItemUpdatedAt.Time)
		}

		// Create map for price
		if price.Id != "" {
			productPrices[price.Id] = price
			productPriceItems[price.Id] = map[string]*v1.PriceItem{}
			productPriceCurrencies[price.Id] = map[string]*v1.PriceCurrency{}
		}

		if productItem.Id != "" {
			productItems[productItem.Id] = productItem
		}

		if productCurrency.Id != "" {
			productCurrencies[productCurrency.Id] = productCurrency
		}

		if priceItem.Id != "" {
			priceItem.Item = priceItemItem
			productPriceItems[price.Id][priceItem.Id] = priceItem
		}

		if priceCurrency.Id != "" {
			priceCurrency.Currency = priceCurrencyCurrency
			productPriceCurrencies[price.Id][priceCurrency.Id] = priceCurrency
		}
	}

	// Convert item map into item slice
	items := []*v1.ProductItem{}
	for _, value := range productItems {
		items = append(items, value)
	}

	// Convert currency map into currency slice
	currencies := []*v1.ProductCurrency{}
	for _, value := range productCurrencies {
		currencies = append(currencies, value)
	}

	prices := []*v1.Price{}
	for _, value := range productPrices {
		prices = append(prices, value)
	}

	for _, price := range prices {
		priceItems := []*v1.PriceItem{}

		for _, item := range productPriceItems[price.Id] {
			priceItems = append(priceItems, item)
		}

		price.Items = priceItems
	}

	for _, price := range prices {
		priceCurrencies := []*v1.PriceCurrency{}

		for _, currency := range productPriceCurrencies[price.Id] {
			priceCurrencies = append(priceCurrencies, currency)
		}

		price.Currencies = priceCurrencies
	}

	product := &v1.Product{
		Id:         res.ProductID,
		Name:       res.ProductName,
		Items:      items,
		Currencies: currencies,
		Prices:     prices,
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
