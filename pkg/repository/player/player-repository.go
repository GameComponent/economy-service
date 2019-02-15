package playerrepository

import (
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
func (r *PlayerRepository) Get(id string) (*v1.Player, error) {
	rows, err := r.db.Query(
		`SELECT id, name
    FROM storage 
    WHERE player_id = $1`,
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
