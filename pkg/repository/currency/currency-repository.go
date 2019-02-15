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
func (r *CurrencyRepository) Create(ctx context.Context, name string) (*v1.Currency, error) {
	// Add item to the databased return the generated UUID
	lastInsertUUID := ""
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO currency(name) VALUES ($1) RETURNING id`,
		name,
	).Scan(&lastInsertUUID)

	if err != nil {
		return nil, err
	}

	// Generate the object based on the generated id and the requested name
	currency := &v1.Currency{
		Id:   lastInsertUUID,
		Name: name,
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
