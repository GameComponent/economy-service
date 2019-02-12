package v1

import (

	// "errors"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// "github.com/golang/protobuf/ptypes/struct"
	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	itemrepository "github.com/GameComponent/economy-service/pkg/repository/item"
	playerrepository "github.com/GameComponent/economy-service/pkg/repository/player"
)

const (
	apiVersion = "v1"
)

// toDoServiceServer is implementation of v1.ToDoServiceServer proto interface
type economyServiceServer struct {
	db               *sql.DB
	itemRepository   *itemrepository.ItemRepository
	playerRepository *playerrepository.PlayerRepository
}

// NewEconomyServiceServer creates ToDo service
func NewEconomyServiceServer(db *sql.DB, itemRepository *itemrepository.ItemRepository, playerRepository *playerrepository.PlayerRepository) v1.EconomyServiceServer {
	return &economyServiceServer{
		db,
		itemRepository,
		playerRepository,
	}
}

// checkAPI checks if the API version requested by client is supported by server
func (s *economyServiceServer) checkAPI(api string) error {
	// API version is "" means use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}
