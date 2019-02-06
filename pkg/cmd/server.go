package cmd

import (
  "context"
  "flag"
  "fmt"
  "log"

  "github.com/GameComponent/economy-service/pkg/protocol/grpc"
  v1 "github.com/GameComponent/economy-service/pkg/service/v1"
  database "github.com/GameComponent/economy-service/pkg/database"
)

const (
  databaseHost      = "localhost"
  databasePort      = "26257"
  databaseUser      = "root"
  databasePassword  = ""
  databaseName      = "economy"
  databaseSsl       = "disable"
)

type Config struct {
  GRPCPort string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
  ctx := context.Background()

  // Setup database & migrate
  _, err := database.Init(
    databaseHost,
    databasePort,
    databaseUser,
    databasePassword,
    databaseName,
    databaseSsl,
  );
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

  // Start the application
  v1API := v1.NewEconomyServiceServer(db)
  return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}

