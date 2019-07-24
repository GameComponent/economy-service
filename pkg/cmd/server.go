package cmd

import (
	"context"
	"flag"
	"fmt"
	"log"

	database "github.com/GameComponent/economy-service/pkg/database"
	"github.com/GameComponent/economy-service/pkg/logger"
	"github.com/GameComponent/economy-service/pkg/protocol/grpc"
	"github.com/GameComponent/economy-service/pkg/protocol/rest"
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
)

// Config for the server
type Config struct {
	GRPCPort         string `mapstructure:"grpc_port"`
	HTTPPort         string `mapstructure:"http_port"`
	DatabaseHost     string `mapstructure:"db_host"`
	DatabasePort     string `mapstructure:"db_port"`
	DatabaseUser     string `mapstructure:"db_user"`
	DatabasePassword string `mapstructure:"db_password"`
	DatabaseName     string `mapstructure:"db_name"`
	DatabaseSsl      string `mapstructure:"db_ssl"`
	LogLevel         int    `mapstructure:"log_level"`
	LogTimeFormat    string `mapstructure:"log_time_format"`
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

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

	// Set potential config locations
	v.SetConfigName("config")
	v.AddConfigPath("/etc/economy-service/")
	v.AddConfigPath("$HOME/.economy-service")
	v.AddConfigPath(".")

	// Load from config files
	err := v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found")
		} else {
			return fmt.Errorf("Unable to read in config")
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

	// Add flags to Viper
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	v.BindPFlags(pflag.CommandLine)

	// Unmarhsal the viper config to a struct
	var cfg Config
	err = v.Unmarshal(&cfg)
	if err != nil {
		return fmt.Errorf("Unable to unmarshal config")
	}

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	// Setup the logger
	if err := logger.Init(cfg.LogLevel, cfg.LogTimeFormat); err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
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
		log.Fatal("Could create a database connection")
	}

	// Setup the repositories
	itemRepository := itemrepository.NewItemRepository(db)
	playerRepository := playerrepository.NewPlayerRepository(db)
	currencyRepository := currencyrepository.NewCurrencyRepository(db)
	storageRepository := storagerepository.NewStorageRepository(db)
	configRepository := configrepository.NewConfigRepository(db)
	accountRepository := accountrepository.NewAccountRepository(db)
	shopRepository := shoprepository.NewShopRepository(db)
	productRepository := productrepository.NewProductRepository(db)
	priceRepository := pricerepository.NewPriceRepository(db)

	// Create the config
	config := v1.Config{
		DB:                 db,
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
		_ = rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
	}()

	// Start the GRCP server
	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
