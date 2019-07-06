package shoprepository

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	repository "github.com/GameComponent/economy-service/pkg/repository"
	"github.com/golang/protobuf/ptypes"
)

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

// ShopRepository struct
type ShopRepository struct {
	db *sql.DB
}

// NewShopRepository constructor
func NewShopRepository(db *sql.DB) repository.ShopRepository {
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
			FROM shop
      LEFT JOIN shop_product ON (shop_product.shop_id = shop.id)
      LEFT JOIN product ON (product.id = shop_product.product_id)
      LEFT JOIN product_item ON (product_item.product_id = product.id)
      LEFT JOIN item ON (item.id = product_item.item_id)
      LEFT JOIN product_currency ON (product_currency.product_id = product.id)
      LEFT JOIN currency ON (currency.id = product_currency.currency_id)
      LEFT JOIN price ON (price.product_id = product.id)
      LEFT JOIN price_currency ON (price_currency.price_id = price.id)
      LEFT JOIN price_item ON (price_item.price_id = price.id)
      LEFT JOIN currency price_currency_currency ON (price_currency_currency.id = price_currency.currency_id)
      LEFT JOIN item price_item_item ON (price_item_item.id = price_item.item_id)
			WHERE shop.id = $1
		`,
		shopID,
	)

	if err != nil {
		return nil, err
	}

	type row struct {
		ShopID                            string
		ShopName                          string
		ShopCreatedAt                     time.Time
		ShopUpdatedAt                     time.Time
		ProductID                         sql.NullString
		ProductName                       sql.NullString
		ProductCreatedAt                  NullTime
		ProductUpdatedAt                  NullTime
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

	// productItems := map[string]map[string]*v1.ProductItem{}
	shopProducts := map[string]*v1.Product{}
	shopProductItems := map[string]map[string]*v1.ProductItem{}
	shopProductCurrencies := map[string]map[string]*v1.ProductCurrency{}
	shopProductPrices := map[string]map[string]*v1.Price{}
	shopProductPriceItems := map[string]map[string]map[string]*v1.PriceItem{}
	shopProductPriceCurrencies := map[string]map[string]map[string]*v1.PriceCurrency{}

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
			item.Stackable = res.ItemStackable.Bool
			item.StackMaxAmount = res.ItemStackMaxAmount.Int64
			item.StackBalancingMethod = v1.StackBalancingMethod(res.ItemStackBalancingMethod.Int64)
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

		// Extract the Currency
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

		// Create the products
		if res.ProductID.Valid && shopProducts[res.ProductID.String] == nil {
			product := &v1.Product{
				Id:   res.ProductID.String,
				Name: res.ProductName.String,
			}

			product.CreatedAt, _ = ptypes.TimestampProto(res.ProductCreatedAt.Time)
			product.UpdatedAt, _ = ptypes.TimestampProto(res.ProductUpdatedAt.Time)

			shopProducts[res.ProductID.String] = product
			shopProductItems[res.ProductID.String] = map[string]*v1.ProductItem{}
			shopProductCurrencies[res.ProductID.String] = map[string]*v1.ProductCurrency{}

			shopProductPrices[res.ProductID.String] = map[string]*v1.Price{}
			shopProductPriceItems[res.ProductID.String] = map[string]map[string]*v1.PriceItem{}
			shopProductPriceCurrencies[res.ProductID.String] = map[string]map[string]*v1.PriceCurrency{}
		}

		// Create map for Price
		if price.Id != "" {
			shopProductPrices[res.ProductID.String][price.Id] = price
			shopProductPriceItems[res.ProductID.String][price.Id] = map[string]*v1.PriceItem{}
			shopProductPriceCurrencies[res.ProductID.String][price.Id] = map[string]*v1.PriceCurrency{}
		}

		if priceItem.Id != "" {
			priceItem.Item = priceItemItem
			shopProductPriceItems[res.ProductID.String][price.Id][priceItem.Id] = priceItem
		}

		if priceCurrency.Id != "" {
			priceCurrency.Currency = priceCurrencyCurrency
			shopProductPriceCurrencies[res.ProductID.String][price.Id][priceCurrency.Id] = priceCurrency
		}

		// Add object to the productItems if it is set
		if res.ProductItemID.Valid && res.ItemID.Valid {
			if shopProductItems[res.ProductID.String] == nil {
				shopProductItems[res.ProductID.String] = map[string]*v1.ProductItem{}
			}
			shopProductItems[res.ProductID.String][productItem.Id] = productItem
		}

		// Add object to the productCurrencies if it is set
		if res.ProductCurrencyID.Valid && res.CurrencyID.Valid {
			if shopProductCurrencies[res.ProductID.String] == nil {
				shopProductCurrencies[res.ProductID.String] = map[string]*v1.ProductCurrency{}
			}
			shopProductCurrencies[res.ProductID.String][productItem.Id] = productCurrency
		}
	}

	// Convert item map into item slice
	products := []*v1.Product{}
	for shopProductKey, shopProductValue := range shopProducts {
		// Add the Items to the Product
		productItems := []*v1.ProductItem{}
		for _, productItemValue := range shopProductItems[shopProductKey] {
			productItems = append(productItems, productItemValue)
		}
		shopProductValue.Items = productItems

		// Add the Currencies to the Product
		productCurrencies := []*v1.ProductCurrency{}
		for _, productCurrencyValue := range shopProductCurrencies[shopProductKey] {
			productCurrencies = append(productCurrencies, productCurrencyValue)
		}
		shopProductValue.Currencies = productCurrencies

		// Add the Prices to the Product
		productPrices := []*v1.Price{}
		for _, price := range shopProductPrices[shopProductKey] {
			// Add the Items to the Price
			priceItems := []*v1.PriceItem{}
			for _, productPriceItemValue := range shopProductPriceItems[shopProductKey][price.Id] {
				priceItems = append(priceItems, productPriceItemValue)
			}
			price.Items = priceItems

			// Add the Currencies to the Price
			priceCurrencies := []*v1.PriceCurrency{}
			for _, currency := range shopProductPriceCurrencies[shopProductKey][price.Id] {
				priceCurrencies = append(priceCurrencies, currency)
			}
			price.Currencies = priceCurrencies

			productPrices = append(productPrices, price)
		}
		shopProductValue.Prices = productPrices

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
func (r *ShopRepository) List(ctx context.Context, limit int32, offset int32) ([]*v1.Shop, int32, error) {
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
