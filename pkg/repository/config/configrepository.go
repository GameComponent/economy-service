package configrepository

import (
	"context"
	"database/sql"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	repository "github.com/GameComponent/economy-service/pkg/repository"
	"go.uber.org/zap"
)

// ConfigRepository struct
type ConfigRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewConfigRepository constructor
func NewConfigRepository(db *sql.DB, logger *zap.Logger) repository.ConfigRepository {
	return &ConfigRepository{
		db:     db,
		logger: logger,
	}
}

// Get a player
func (r *ConfigRepository) Get(ctx context.Context, key string) (*v1.Config, error) {
	config := &v1.Config{
		Key:   key,
	}

	err := r.db.QueryRowContext(
		ctx,
		`SELECT value FROM config WHERE key = $1`,
		key,
	).Scan(&config.Value)

	if err != nil {
		return nil, err
	}

	return config, nil
}

// Set a new config
func (r *ConfigRepository) Set(ctx context.Context, key string, value string) (*v1.Config, error) {
	_, err := r.db.ExecContext(
		ctx,
		`
			INSERT INTO config(key, value)
			VALUES ($1, $2)
			ON CONFLICT(key)
			DO UPDATE
			SET value = excluded.value
		`,
		key,
		value,
	)

	if err != nil {
		return nil, err
	}

	return &v1.Config{
		Key:   key,
		Value: value,
	}, nil
}

// List all configs
func (r *ConfigRepository) List(ctx context.Context, limit int32, offset int32) ([]*v1.Config, int32, error) {
	// Query configs from the database
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT 
				key,
				value,
				(SELECT COUNT(*) FROM config) AS total_size
			FROM config
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

	// Unwrap rows into configs
	configs := []*v1.Config{}
	totalSize := int32(0)

	for rows.Next() {
		config := v1.Config{}

		err := rows.Scan(
			&config.Key,
			&config.Value,
			&totalSize,
		)
		if err != nil {
			return nil, 0, err
		}

		configs = append(configs, &config)
	}

	return configs, totalSize, nil
}
