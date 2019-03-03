package playerrepository

import (
	"context"
	"database/sql"
	"time"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	"github.com/golang/protobuf/ptypes"
)

// PlayerRepository struct
type PlayerRepository struct {
	db *sql.DB
}

// NewPlayerRepository constructor
func NewPlayerRepository(db *sql.DB) *PlayerRepository {
	return &PlayerRepository{
		db: db,
	}
}

// Get a player
func (r *PlayerRepository) Get(ctx context.Context, id string) (*v1.Player, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT id, name, created_at, updated_at
			FROM storage 
			WHERE player_id = $1
		`,
		id,
	)

	if err != nil {
		return nil, err
	}

	storages := []*v1.Storage{}

	for rows.Next() {
		storage := &v1.Storage{}
		createdAt := time.Time{}
		updatedAt := time.Time{}

		err = rows.Scan(
			&storage.Id,
			&storage.Name,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Convert created_at to timestamp
		storage.CreatedAt, _ = ptypes.TimestampProto(createdAt)
		storage.UpdatedAt, _ = ptypes.TimestampProto(updatedAt)

		storages = append(storages, storage)
	}

	return &v1.Player{
		Id:       id,
		Storages: storages,
	}, nil
}

// List all player
func (r *PlayerRepository) List(
	ctx context.Context,
	limit int32,
	offset int32,
) (
	[]*v1.Player,
	int32,
	error,
) {
	// Query items from the database
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT DISTINCT
				player_id,
				(SELECT COUNT(DISTINCT player_id) FROM storage) AS total_size
			FROM storage
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
			&totalSize,
		)
		if err != nil {
			return nil, 0, err
		}

		players = append(players, &player)
	}

	return players, totalSize, nil
}
