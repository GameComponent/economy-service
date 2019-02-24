package currencyrepository

import (
	"context"
	"database/sql"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
)

// CurrencyRepository struct
type CurrencyRepository struct {
	db *sql.DB
}

// NewCurrencyRepository constructor
func NewCurrencyRepository(db *sql.DB) *CurrencyRepository {
	return &CurrencyRepository{
		db: db,
	}
}

// Create a currency
func (r *CurrencyRepository) Create(
	ctx context.Context,
	name string,
	shortName string,
	symbol string,
) (*v1.Currency, error) {
	// Add item to the databased return the generated UUID
	lastInsertUUID := ""
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO currency(name, short_name, symbol) VALUES ($1, $2, $3) RETURNING id`,
		name,
		shortName,
		symbol,
	).Scan(
		&lastInsertUUID,
	)

	if err != nil {
		return nil, err
	}

	// Generate the object based on the generated id and the requested name
	currency := &v1.Currency{
		Id:        lastInsertUUID,
		Name:      name,
		ShortName: shortName,
		Symbol:    symbol,
	}

	return currency, nil
}

// Get a currency
func (r *CurrencyRepository) Get(ctx context.Context, currencyID string) (*v1.Currency, error) {
	currency := &v1.Currency{}

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, name FROM currency WHERE id = $1`,
		currencyID,
	).Scan(
		&currency.Id,
		&currency.Name,
	)

	if err != nil {
		return nil, err
	}

	return currency, nil
}

// List all currenciees
// func (r *CurrencyRepository) List(ctx context.Context) ([]*v1.Currency, error) {
// 	currency := &v1.Currency{}

// 	err := r.db.QueryRowContext(
// 		ctx,
// 		`SELECT id, name FROM currency`,
// 	).Scan(
// 		&currency.Id,
// 		&currency.Name,
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return currency, nil
// }
