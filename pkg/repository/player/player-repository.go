package playerrepository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	repository "github.com/GameComponent/economy-service/pkg/repository"
	jsonpb "github.com/golang/protobuf/jsonpb"
	_struct "github.com/golang/protobuf/ptypes/struct"
)

// PlayerRepository struct
type PlayerRepository struct {
	db *sql.DB
}

// NewPlayerRepository constructor
func NewPlayerRepository(db *sql.DB) repository.PlayerRepository {
	return &PlayerRepository{
		db: db,
	}
}

// Create a new player
func (r *PlayerRepository) Create(ctx context.Context, playerID string, name string, metadata *_struct.Struct) (*v1.Player, error) {
	// Parse struct to JSON string
	jsonMetadata := "{}"
	if metadata != nil {
		var err error
		marshaler := jsonpb.Marshaler{}
		jsonMetadata, err = marshaler.MarshalToString(metadata)
		if err != nil {
			return nil, err
		}
	}

	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO player(id, name, metadata) VALUES ($1, $2, $3)`,
		playerID,
		name,
		jsonMetadata,
	)

	if err != nil {
		return nil, err
	}

	return r.Get(ctx, playerID)
}

// Update a player
func (r *PlayerRepository) Update(ctx context.Context, playerID string, name string, metadata *_struct.Struct) (*v1.Player, error) {
	index := 1
	queries := []string{}
	arguments := []interface{}{}

	// Add name to the query
	if name != "" {
		queries = append(queries, fmt.Sprintf("name = $%v", index))
		arguments = append(arguments, name)
		index++
	}

	// Add metadata to the query
	if metadata != nil {
		// Parse the metadata to a JSON string
		jsonMetadata := "{}"
		var err error
		marshaler := jsonpb.Marshaler{}
		jsonMetadata, err = marshaler.MarshalToString(metadata)
		if err != nil {
			return nil, err
		}

		queries = append(queries, fmt.Sprintf("metadata = $%v", index))
		arguments = append(arguments, jsonMetadata)
		index++
	}

	if index <= 1 {
		return nil, fmt.Errorf("no arguments given")
	}

	// Update the player
	arguments = append(arguments, playerID)
	query := fmt.Sprintf("UPDATE player SET %v WHERE id =$%v", strings.Join(queries, ", "), index)
	_, err := r.db.ExecContext(
		ctx,
		query,
		arguments...,
	)

	if err != nil {
		return nil, err
	}

	return r.Get(ctx, playerID)
}

// Get a player
func (r *PlayerRepository) Get(ctx context.Context, playerID string) (*v1.Player, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT 
				player.id AS playerId,
				player.name AS playerName,
				player.metadata AS playerMetadata,
				storage.id as storageId,
				storage.name as storageName
			FROM player
			LEFT JOIN storage ON (player.id = storage.player_id)
			WHERE player.id = $1
		`,
		playerID,
	)
	if err != nil {
		return nil, err
	}

	type row struct {
		PlayerID       string
		PlayerName     string
		PlayerMetadata string
		StorageID      sql.NullString
		StorageName    sql.NullString
	}

	storages := []*v1.Storage{}

	var res row
	for rows.Next() {
		err = rows.Scan(
			&res.PlayerID,
			&res.PlayerName,
			&res.PlayerMetadata,
			&res.StorageID,
			&res.StorageName,
		)
		if err != nil {
			return nil, err
		}

		if res.StorageID.Valid {
			storage := v1.Storage{
				Id:   res.StorageID.String,
				Name: res.StorageName.String,
			}
			storages = append(storages, &storage)
		}
	}

	// Check if there is atleast 1 row found
	if res.PlayerID == "" {
		return nil, fmt.Errorf("Player not found")
	}

	// Create the player struct
	player := &v1.Player{
		Id:       res.PlayerID,
		Name:     res.PlayerName,
		Storages: storages,
	}

	// Convert metadata json to a proto Struct
	if res.PlayerMetadata != "" {
		metadataStruct := _struct.Struct{}
		stringReader := strings.NewReader(res.PlayerMetadata)
		unmarshaler := jsonpb.Unmarshaler{}
		err = unmarshaler.Unmarshal(stringReader, &metadataStruct)
		if err != nil {
			return nil, err
		}

		player.Metadata = &metadataStruct
	}

	return player, nil
}

// List all player
func (r *PlayerRepository) List(ctx context.Context, limit int32, offset int32) ([]*v1.Player, int32, error) {
	// Query items from the database
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT
				id,
				name,
				(SELECT COUNT(id) FROM player) AS total_size
			FROM player
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

	// Unwrap rows into items
	players := []*v1.Player{}
	totalSize := int32(1)

	for rows.Next() {
		player := v1.Player{}

		err := rows.Scan(
			&player.Id,
			&player.Name,
			&totalSize,
		)
		if err != nil {
			return nil, 0, err
		}

		players = append(players, &player)
	}

	return players, totalSize, nil
}

// Search player
func (r *PlayerRepository) Search(ctx context.Context, query string, limit int32, offset int32) ([]*v1.Player, int32, error) {
	// Query items from the database
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT
				id,
				name,
				(SELECT COUNT(id) FROM player WHERE name ~* $1) AS total_size
			FROM player
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

	// Unwrap rows into items
	players := []*v1.Player{}
	totalSize := int32(0)

	for rows.Next() {
		player := v1.Player{}

		err := rows.Scan(
			&player.Id,
			&player.Name,
			&totalSize,
		)
		if err != nil {
			return nil, 0, err
		}

		players = append(players, &player)
	}

	return players, totalSize, nil
}
