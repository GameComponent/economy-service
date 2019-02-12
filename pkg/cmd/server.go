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
	itemrepository "github.com/GameComponent/economy-service/pkg/repository/item"
	playerrepository "github.com/GameComponent/economy-service/pkg/repository/player"
	v1 "github.com/GameComponent/economy-service/pkg/service/v1"
)

const (
	databaseHost     = "localhost"
	databasePort     = "26257"
	databaseUser     = "root"
	databasePassword = ""
	databaseName     = "economy"
	databaseSsl      = "disable"
	HTTPPort         = "8888"
	LogLevel         = 0
	LogTimeFormat    = ""
)

// Config for the server
type Config struct {
	GRPCPort string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// Setup the logger
	if err := logger.Init(LogLevel, LogTimeFormat); err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
	}

	// Setup database & migrate
	_, err := database.Init(
		databaseHost,
		databasePort,
		databaseUser,
		databasePassword,
		databaseName,
		databaseSsl,
	)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(
		databaseHost,
		databasePort,
		databaseUser,
		databasePassword,
		databaseName,
		databaseSsl,
	)
	defer db.Close()

	if err != nil {
		log.Fatal("Could create a database connection")
	}

	// Get configuration
	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc_port", "", "gRPC port to bind")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	// Setup the repositories
	itemRepository := itemrepository.NewItemRepository(db)
	playerRepository := playerrepository.NewPlayerRepository(db)

	// Start the service
	v1API := v1.NewEconomyServiceServer(db, itemRepository, playerRepository)

	// Start the REST server
	go func() {
		_ = rest.RunServer(ctx, cfg.GRPCPort, HTTPPort)
	}()

	// Start the GRCP server
	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
