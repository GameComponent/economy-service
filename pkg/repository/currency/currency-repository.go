package currencyrepository

import (
	"context"
	"database/sql"
	"time"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	repository "github.com/GameComponent/economy-service/pkg/repository"
	"github.com/golang/protobuf/ptypes"
)

// CurrencyRepository struct
type CurrencyRepository struct {
	db *sql.DB
}

// NewCurrencyRepository constructor
func NewCurrencyRepository(db *sql.DB) repository.CurrencyRepository {
	return &CurrencyRepository{
		db: db,
	}
}

// Create a currency
func (r *CurrencyRepository) Create(ctx context.Context, name string, shortName string, symbol string) (*v1.Currency, error) {
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

// Update a currency
func (r *CurrencyRepository) Update(ctx context.Context, currencyID string, name string, shortName string, symbol string) (*v1.Currency, error) {
	_, err := r.db.ExecContext(
		ctx,
		`
			UPDATE currency
			SET name = $1, short_name = $2, symbol = $3
			WHERE id = $4
		`,
		name,
		shortName,
		symbol,
		currencyID,
	)
	if err != nil {
		return nil, err
	}

	return r.Get(ctx, currencyID)
}

// Get a currency
func (r *CurrencyRepository) Get(ctx context.Context, currencyID string) (*v1.Currency, error) {
	currency := &v1.Currency{}

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, short_name, symbol FROM currency WHERE id = $1`,
		currencyID,
	).Scan(
		&currency.Id,
		&currency.Name,
		&currency.ShortName,
		&currency.Symbol,
	)

	if err != nil {
		return nil, err
	}

	return currency, nil
}

// List all currenciees
func (r *CurrencyRepository) List(ctx context.Context, limit int32, offset int32) ([]*v1.Currency, int32, error) {
	// Query items from the database
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT 
				id,
				name,
				short_name,
				symbol,
				created_at,
				updated_at,
				(SELECT COUNT(*) FROM currency) AS total_size
			FROM currency
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

	// Unwrap rows into currency
	currencies := []*v1.Currency{}
	totalSize := int32(0)

	for rows.Next() {
		currency := v1.Currency{}
		createdAt := time.Time{}
		updatedAt := time.Time{}

		err := rows.Scan(
			&currency.Id,
			&currency.Name,
			&currency.ShortName,
			&currency.Symbol,
			&createdAt,
			&updatedAt,
			&totalSize,
		)
		if err != nil {
			return nil, 0, err
		}

		// Convert created_at to timestamp
		currency.CreatedAt, _ = ptypes.TimestampProto(createdAt)
		currency.UpdatedAt, _ = ptypes.TimestampProto(updatedAt)

		currencies = append(currencies, &currency)
	}

	return currencies, totalSize, nil
}
