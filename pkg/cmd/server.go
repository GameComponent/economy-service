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
	configrepository "github.com/GameComponent/economy-service/pkg/repository/config"
	currencyrepository "github.com/GameComponent/economy-service/pkg/repository/currency"
	itemrepository "github.com/GameComponent/economy-service/pkg/repository/item"
	playerrepository "github.com/GameComponent/economy-service/pkg/repository/player"
	storagerepository "github.com/GameComponent/economy-service/pkg/repository/storage"
	v1 "github.com/GameComponent/economy-service/pkg/service/v1"
)

// Config for the server
type Config struct {
	GRPCPort         string
	HTTPPort         string
	DatabaseHost     string
	DatabasePort     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	DatabaseSsl      string
	LogLevel         int
	LogTimeFormat    string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// Get configuration
	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc_port", "3000", "gRPC port to bind")
	flag.StringVar(&cfg.HTTPPort, "http_port", "8080", "http port to bind")
	flag.StringVar(&cfg.DatabaseHost, "db_host", "127.0.0.1", "host of the database")
	flag.StringVar(&cfg.DatabasePort, "db_port", "26257", "port of the database")
	flag.StringVar(&cfg.DatabaseUser, "db_user", "root", "user of the database")
	flag.StringVar(&cfg.DatabasePassword, "db_password", "", "password of the database")
	flag.StringVar(&cfg.DatabaseName, "db_name", "economy", "name of the database")
	flag.StringVar(&cfg.DatabaseSsl, "db_ssl", "disable", "ssl settings of the database")
	flag.IntVar(&cfg.LogLevel, "log_level", 0, "level of the logger")
	flag.StringVar(&cfg.LogTimeFormat, "log_time_format", "", "time format of the logger")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	// Setup the logger
	if err := logger.Init(cfg.LogLevel, cfg.LogTimeFormat); err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
	}

	// Setup database & migrate
	_, err := database.Init(
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

	// Start the service
	v1API := v1.NewEconomyServiceServer(
		db,
		itemRepository,
		playerRepository,
		currencyRepository,
		storageRepository,
		configRepository,
	)

	// Start the REST server
	go func() {
		_ = rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
	}()

	// Start the GRCP server
	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
