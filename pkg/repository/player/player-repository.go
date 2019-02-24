package playerrepository

import (
	"context"
	"database/sql"
	"log"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
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
			SELECT id, name
			FROM storage 
			WHERE player_id = $1
		`,
		id,
	)

	if err != nil {
		return nil, err
	}

	storageItems := []*v1.StorageBase{}

	for rows.Next() {
		storage := &v1.StorageBase{}

		err = rows.Scan(
			&storage.Id,
			&storage.Name,
		)
		if err != nil {
			log.Fatal(err)
		}

		storageItems = append(storageItems, storage)
	}

	return &v1.Player{
		Id:       id,
		Storages: storageItems,
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
			SELECT 
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
