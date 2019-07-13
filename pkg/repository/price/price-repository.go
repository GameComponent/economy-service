package pricerepository

import (
	"context"
	"database/sql"
	"time"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	repository "github.com/GameComponent/economy-service/pkg/repository"
	"github.com/golang/protobuf/ptypes"
)

// PriceRepository struct
type PriceRepository struct {
	db *sql.DB
}

// NewPriceRepository constructor
func NewPriceRepository(db *sql.DB) repository.PriceRepository {
	return &PriceRepository{
		db: db,
	}
}

// Get a price
func (r *PriceRepository) Get(ctx context.Context, priceID string) (*v1.Price, error) {
	createdAt := time.Time{}
	updatedAt := time.Time{}

	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT
				price.id AS priceId,
				price.created_at AS priceCreatedAt,
				price.updated_at AS priceUpdatedAt,
				price_currency.id AS priceCurrencyId,
				price_currency.amount as priceCurrencyAmount,
				price_item.id as priceItemId,
				price_item.amount as priceItemAmount,
				currency.id AS currencyId,
				currency.name AS currencyName,
				currency.short_name AS currencyShortName,
				currency.symbol AS currencySymbol,
				item.id as itemId,
				item.name as itemName,
				item.stackable as itemStackable,
				item.stack_max_amount as itemStackMaxAmount,
				item.stack_balancing_method as itemStackBalancingMethod,
				item.metadata as itemData
			FROM price
			LEFT JOIN price_item ON (price.id = price_item.price_id)
			LEFT JOIN item ON (price_item.item_id = item.id)
			LEFT JOIN price_currency ON (price.id = price_currency.price_id)
			LEFT JOIN currency ON (currency.id = price_currency.currency_id)
			WHERE price.id = $1
		`,
		priceID,
	)
	if err != nil {
		return nil, err
	}

	type row struct {
		PriceID                  string
		PriceCreatedAt           string
		PriceUpdatedAt           string
		PriceCurrencyID          sql.NullString
		PriceCurrencyAmount      sql.NullInt64
		PriceItemID              sql.NullString
		PriceItemAmount          sql.NullInt64
		CurrencyID               sql.NullString
		CurrencyName             sql.NullString
		CurrencyShortName        sql.NullString
		CurrencySymbol           sql.NullString
		ItemID                   sql.NullString
		ItemName                 sql.NullString
		ItemStackable            sql.NullBool
		ItemStackMaxAmount       sql.NullInt64
		ItemStackBalancingMethod sql.NullInt64
		ItemData                 sql.NullString
	}

	priceCurrencies := map[string]*v1.PriceCurrency{}
	priceItems := map[string]*v1.PriceItem{}

	var res row
	for rows.Next() {
		err = rows.Scan(
			&res.PriceID,
			&createdAt,
			&updatedAt,
			&res.PriceCurrencyID,
			&res.PriceCurrencyAmount,
			&res.PriceItemID,
			&res.PriceItemAmount,
			&res.CurrencyID,
			&res.CurrencyName,
			&res.CurrencyShortName,
			&res.CurrencySymbol,
			&res.ItemID,
			&res.ItemName,
			&res.ItemStackable,
			&res.ItemStackMaxAmount,
			&res.ItemStackBalancingMethod,
			&res.ItemData,
		)

		if err != nil {
			return nil, err
		}

		// Extract the Currency
		currency := &v1.Currency{}
		if res.CurrencyID.Valid {
			currency.Id = res.CurrencyID.String
			currency.Name = res.CurrencyName.String
			currency.ShortName = res.CurrencyShortName.String
			currency.Symbol = res.CurrencySymbol.String
		}

		// Extract the PriceCurrency
		priceCurrency := &v1.PriceCurrency{}
		if res.PriceCurrencyID.Valid {
			priceCurrency.Id = res.PriceCurrencyID.String
			priceCurrency.Amount = res.PriceCurrencyAmount.Int64
			priceCurrency.Currency = currency
		}

		// Add object currency to the price if it is set
		if priceCurrency.Id != "" {
			priceCurrencies[priceCurrency.Id] = priceCurrency
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

		// Extract the PriceItem
		priceItem := &v1.PriceItem{}
		if res.PriceItemID.Valid {
			priceItem.Id = res.PriceItemID.String
			priceItem.Item = item

			// Only show the amount if the item is stackable
			if res.ItemStackable.Bool {
				priceItem.Amount = res.PriceItemAmount.Int64
			}
		}

		// Add object to item the price if it is set
		if priceItem.Id != "" {
			priceItems[priceItem.Id] = priceItem
		}
	}

	// Convert item map into item slice
	currencies := []*v1.PriceCurrency{}
	for _, value := range priceCurrencies {
		currencies = append(currencies, value)
	}

	// Convert item map into item slice
	items := []*v1.PriceItem{}
	for _, value := range priceItems {
		items = append(items, value)
	}

	// Get the price
	price := &v1.Price{
		Id:         res.PriceID,
		Currencies: currencies,
		Items:      items,
	}
	price.CreatedAt, _ = ptypes.TimestampProto(createdAt)
	price.UpdatedAt, _ = ptypes.TimestampProto(updatedAt)

	return price, nil
}

// Create a new price
func (r *PriceRepository) Create(ctx context.Context, productID string) (*v1.Price, error) {
	lastInsertUUID := ""
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO price(product_id) VALUES($1) RETURNING id`,
		productID,
	).Scan(&lastInsertUUID)

	if err != nil {
		return nil, err
	}

	return &v1.Price{
		Id: lastInsertUUID,
	}, nil
}

// Delete a Price
func (r *PriceRepository) Delete(ctx context.Context, priceID string) (bool, error) {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM price WHERE id = $1`,
		priceID,
	)

	if err != nil {
		return false, err
	}

	return true, nil
}

// AttachPriceCurrency attaches a currency to a price
func (r *PriceRepository) AttachPriceCurrency(ctx context.Context, priceID string, currencyID string, amount int64) (*v1.Price, error) {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO price_currency(price_id, currency_id, amount) VALUES ($1, $2, $3)`,
		priceID,
		currencyID,
		amount,
	)

	if err != nil {
		return nil, err
	}

	return r.Get(ctx, priceID)
}

// DetachPriceCurrency detaches a currency from a price
func (r *PriceRepository) DetachPriceCurrency(ctx context.Context, priceCurrencyID string) (*v1.Price, error) {
	priceID := ""
	err := r.db.QueryRowContext(
		ctx,
		`DELETE FROM price_currency WHERE id = $1 returning price_id`,
		priceCurrencyID,
	).Scan(&priceID)

	if err != nil {
		return nil, err
	}

	return r.Get(ctx, priceID)
}

// AttachPriceItem attaches an item to a price
func (r *PriceRepository) AttachPriceItem(ctx context.Context, priceID string, itemID string, amount int64) (*v1.Price, error) {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO price_item(price_id, item_id, amount) VALUES ($1, $2, $3)`,
		priceID,
		itemID,
		amount,
	)

	if err != nil {
		return nil, err
	}

	return r.Get(ctx, priceID)
}

// DetachPriceItem detaches an item from a price
func (r *PriceRepository) DetachPriceItem(ctx context.Context, priceItemID string) (*v1.Price, error) {
	priceID := ""
	err := r.db.QueryRowContext(
		ctx,
		`DELETE FROM price_item WHERE id = $1 returning price_id`,
		priceItemID,
	).Scan(&priceID)

	if err != nil {
		return nil, err
	}

	return r.Get(ctx, priceID)
}
