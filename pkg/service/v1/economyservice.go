package v1

import (
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	config "github.com/GameComponent/economy-service/pkg/config"
	repository "github.com/GameComponent/economy-service/pkg/repository"
	"go.uber.org/zap"
)

const (
	apiVersion = "v1"
)

// Config for the server
type Config struct {
	DB                 *sql.DB
	Logger             *zap.Logger
	Config             *config.Config
	AccountRepository  repository.AccountRepository
	ConfigRepository   repository.ConfigRepository
	CurrencyRepository repository.CurrencyRepository
	ItemRepository     repository.ItemRepository
	PlayerRepository   repository.PlayerRepository
	PriceRepository    repository.PriceRepository
	ProductRepository  repository.ProductRepository
	ShopRepository     repository.ShopRepository
	StorageRepository  repository.StorageRepository
}

// EconomyServiceServer is implementation of v1.EconomyServiceServer proto interface
type EconomyServiceServer struct {
	DB                 *sql.DB
	Logger             *zap.Logger
	Config             *config.Config
	AccountRepository  repository.AccountRepository
	ConfigRepository   repository.ConfigRepository
	CurrencyRepository repository.CurrencyRepository
	ItemRepository     repository.ItemRepository
	PlayerRepository   repository.PlayerRepository
	PriceRepository    repository.PriceRepository
	ProductRepository  repository.ProductRepository
	ShopRepository     repository.ShopRepository
	StorageRepository  repository.StorageRepository
}

// NewEconomyServiceServer creates economy service
func NewEconomyServiceServer(config Config) v1.EconomyServiceServer {
	return &EconomyServiceServer{
		config.DB,
		config.Logger,
		config.Config,
		config.AccountRepository,
		config.ConfigRepository,
		config.CurrencyRepository,
		config.ItemRepository,
		config.PlayerRepository,
		config.PriceRepository,
		config.ProductRepository,
		config.ShopRepository,
		config.StorageRepository,
	}
}

// checkAPI checks if the API version requested by client is supported by server
func (s *EconomyServiceServer) checkAPI(api string) error {
	// API version is "" means use current version of the service
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
		}
	}
	return nil
}
