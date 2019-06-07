package v1

import (

	// "errors"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// "github.com/golang/protobuf/ptypes/struct"
	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	accountrepository "github.com/GameComponent/economy-service/pkg/repository/account"
	configrepository "github.com/GameComponent/economy-service/pkg/repository/config"
	currencyrepository "github.com/GameComponent/economy-service/pkg/repository/currency"
	itemrepository "github.com/GameComponent/economy-service/pkg/repository/item"
	playerrepository "github.com/GameComponent/economy-service/pkg/repository/player"
	pricerepository "github.com/GameComponent/economy-service/pkg/repository/price"
	productrepository "github.com/GameComponent/economy-service/pkg/repository/product"
	shoprepository "github.com/GameComponent/economy-service/pkg/repository/shop"
	storagerepository "github.com/GameComponent/economy-service/pkg/repository/storage"
)

const (
	apiVersion = "v1"
)

// economyServiceServer is implementation of v1.EconomyServiceServer proto interface
type economyServiceServer struct {
	db                 *sql.DB
	itemRepository     *itemrepository.ItemRepository
	playerRepository   *playerrepository.PlayerRepository
	currencyRepository *currencyrepository.CurrencyRepository
	storageRepository  *storagerepository.StorageRepository
	configRepository   *configrepository.ConfigRepository
	accountRepository  *accountrepository.AccountRepository
	shopRepository     *shoprepository.ShopRepository
	productRepository  *productrepository.ProductRepository
	priceRepository    *pricerepository.PriceRepository
}

// NewEconomyServiceServer creates economy service
func NewEconomyServiceServer(
	db *sql.DB,
	itemRepository *itemrepository.ItemRepository,
	playerRepository *playerrepository.PlayerRepository,
	currencyRepository *currencyrepository.CurrencyRepository,
	storageRepository *storagerepository.StorageRepository,
	configRepository *configrepository.ConfigRepository,
	accountRepository *accountrepository.AccountRepository,
	shopRepository *shoprepository.ShopRepository,
	productRepository *productrepository.ProductRepository,
	priceRepository *pricerepository.PriceRepository,
) v1.EconomyServiceServer {
	return &economyServiceServer{
		db,
		itemRepository,
		playerRepository,
		currencyRepository,
		storageRepository,
		configRepository,
		accountRepository,
		shopRepository,
		productRepository,
		priceRepository,
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
