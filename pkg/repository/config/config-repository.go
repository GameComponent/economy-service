package configrepository

import (
	"context"
	"database/sql"
	"strings"

	"github.com/golang/protobuf/jsonpb"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	_struct "github.com/golang/protobuf/ptypes/struct"
)

// ConfigRepository struct
type ConfigRepository struct {
	db *sql.DB
}

// NewConfigRepository constructor
func NewConfigRepository(db *sql.DB) *ConfigRepository {
	return &ConfigRepository{
		db: db,
	}
}

// Get a player
func (r *ConfigRepository) Get(ctx context.Context, key string) (*v1.Config, error) {
	var jsonString string

	err := r.db.QueryRowContext(
		ctx,
		`SELECT value FROM config WHERE key = $1`,
		key,
	).Scan(&jsonString)

	if err != nil {
		return nil, err
	}

	stringReader := strings.NewReader(jsonString)
	valueStruct := _struct.Struct{}
	unmarshaler := jsonpb.Unmarshaler{}

	err = unmarshaler.Unmarshal(stringReader, &valueStruct)
	if err != nil {
		return nil, err
	}

	config := &v1.Config{
		Key:   key,
		Value: &valueStruct,
	}

	return config, nil
}

// Set a new config
func (r *ConfigRepository) Set(ctx context.Context, key string, value *_struct.Struct) (*v1.Config, error) {
	marshaler := jsonpb.Marshaler{}
	jsonValue, err := marshaler.MarshalToString(value)
	if err != nil {
		return nil, err
	}

	_, err = r.db.ExecContext(
		ctx,
		`
			INSERT INTO config(key, value)
			VALUES ($1, $2)
			ON CONFLICT(key)
			DO UPDATE
			SET value = excluded.value
		`,
		key,
		jsonValue,
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
func (r *ConfigRepository) List(
	ctx context.Context,
	limit int32,
	offset int32,
) (
	[]*v1.Config,
	int32,
	error,
) {
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
		var jsonString string

		err := rows.Scan(
			&config.Key,
			&jsonString,
			&totalSize,
		)
		if err != nil {
			return nil, 0, err
		}

		stringReader := strings.NewReader(jsonString)
		unmarshaler := jsonpb.Unmarshaler{}
		valueStruct := _struct.Struct{}
		err = unmarshaler.Unmarshal(stringReader, &valueStruct)
		if err != nil {
			return nil, 0, err
		}

		config.Value = &valueStruct

		configs = append(configs, &config)
	}

	return configs, totalSize, nil
}
