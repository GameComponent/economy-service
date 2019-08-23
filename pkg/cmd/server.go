package cmd

import (
	"context"
	"flag"
	"log"

	config "github.com/GameComponent/economy-service/pkg/config"
	database "github.com/GameComponent/economy-service/pkg/database"
	grpc "github.com/GameComponent/economy-service/pkg/protocol/grpc"
	rest "github.com/GameComponent/economy-service/pkg/protocol/rest"
	accountrepository "github.com/GameComponent/economy-service/pkg/repository/account"
	configrepository "github.com/GameComponent/economy-service/pkg/repository/config"
	currencyrepository "github.com/GameComponent/economy-service/pkg/repository/currency"
	itemrepository "github.com/GameComponent/economy-service/pkg/repository/item"
	playerrepository "github.com/GameComponent/economy-service/pkg/repository/player"
	pricerepository "github.com/GameComponent/economy-service/pkg/repository/price"
	productrepository "github.com/GameComponent/economy-service/pkg/repository/product"
	shoprepository "github.com/GameComponent/economy-service/pkg/repository/shop"
	storagerepository "github.com/GameComponent/economy-service/pkg/repository/storage"
	v1 "github.com/GameComponent/economy-service/pkg/service/v1"
	pflag "github.com/spf13/pflag"
	viper "github.com/spf13/viper"
	"go.uber.org/zap"
)

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// Create the logger
	logger, _ := zap.NewProduction()

	// Set the configuration
	v := viper.New()

	// Load from environment variables
	v.AutomaticEnv()

	// Set the defaults
	v.SetDefault("grpc_port", "3000")
	v.SetDefault("http_port", "8080")
	v.SetDefault("db_host", "127.0.0.1")
	v.SetDefault("db_port", "26257")
	v.SetDefault("db_user", "root")
	v.SetDefault("db_password", "")
	v.SetDefault("db_name", "economy")
	v.SetDefault("db_ssl", "disable")
	v.SetDefault("log_level", "0")
	v.SetDefault("log_time_format", "")
	v.SetDefault("jwt_secret", "my_secret_key")
	v.SetDefault("jwt_expiration", 300)

	// Set potential config locations
	v.SetConfigName("config")
	v.AddConfigPath("/etc/economy-service/")
	v.AddConfigPath("$HOME/.economy-service")
	v.AddConfigPath(".")

	// Load from config files
	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Info("No config file found")
		} else {
			logger.Error("Unable to read in config", zap.Error(err))
			return err
		}
	}

	// Parse the flags
	flag.String("grpc_port", "3000", "gRPC port to bind")
	flag.String("http_port", "8080", "http port to bind")
	flag.String("db_host", "127.0.0.1", "host of the database")
	flag.String("db_port", "26257", "port of the database")
	flag.String("db_user", "root", "user of the database")
	flag.String("db_password", "", "password of the database")
	flag.String("db_name", "economy", "name of the database")
	flag.String("db_ssl", "disable", "ssl settings of the database")
	flag.Int("log_level", 0, "level of the logger")
	flag.String("log_time_format", "", "time format of the logger")
	flag.String("jwt_secret", "my_secret_key", "secret used to sign JWT tokens")
	flag.Int("jwt_expiration", 300, "seconds before the JWT expires")

	// Add flags to Viper
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	v.BindPFlags(pflag.CommandLine)

	// Unmarhsal the viper config to a struct
	var cfg config.Config
	err = v.Unmarshal(&cfg)
	if err != nil {
		logger.Error("Unable to unmarshal config", zap.Error(err))
		return err
	}

	if len(cfg.GRPCPort) == 0 {
		logger.Error("invalid TCP port for gRPC server", zap.String("port", cfg.GRPCPort))
		return err
	}

	// Setup database & migrate
	_, err = database.Init(
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseName,
		cfg.DatabaseSsl,
	)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(
		cfg.DatabaseHost,
		cfg.DatabasePort,
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseName,
		cfg.DatabaseSsl,
	)
	defer db.Close()

	if err != nil {
		logger.Fatal("Could create a database connection", zap.Error(err))
	}

	// Setup the repositories
	itemRepository := itemrepository.NewItemRepository(db, logger)
	playerRepository := playerrepository.NewPlayerRepository(db, logger)
	currencyRepository := currencyrepository.NewCurrencyRepository(db, logger)
	storageRepository := storagerepository.NewStorageRepository(db, logger)
	configRepository := configrepository.NewConfigRepository(db, logger)
	accountRepository := accountrepository.NewAccountRepository(db, logger)
	shopRepository := shoprepository.NewShopRepository(db, logger)
	productRepository := productrepository.NewProductRepository(db, logger)
	priceRepository := pricerepository.NewPriceRepository(db, logger)

	// Create the config
	config := v1.Config{
		DB:                 db,
		Logger:             logger,
		Config:             &cfg,
		ItemRepository:     itemRepository,
		PlayerRepository:   playerRepository,
		CurrencyRepository: currencyRepository,
		StorageRepository:  storageRepository,
		ConfigRepository:   configRepository,
		AccountRepository:  accountRepository,
		ShopRepository:     shopRepository,
		ProductRepository:  productRepository,
		PriceRepository:    priceRepository,
	}

	// Start the service
	v1API := v1.NewEconomyServiceServer(config)

	// Start the REST server
	go func() {
		_ = rest.RunServer(ctx, logger, cfg.GRPCPort, cfg.HTTPPort)
	}()

	// Start the GRCP server
	return grpc.RunServer(ctx, v1API, logger, cfg.GRPCPort)
}
