package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	"github.com/GameComponent/economy-service/pkg/logger"
	"github.com/GameComponent/economy-service/pkg/protocol/grpc/middleware"
)

// RunServer runs gRPC service to publish Economy service
func RunServer(ctx context.Context, v1API v1.EconomyServiceServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// gRPC server statup options
	opts := []grpc.ServerOption{}

	// Add middleware
	newOpts := middleware.AddLogging(logger.Log, opts)

	// Register service
	server := grpc.NewServer()
	v1.RegisterEconomyServiceServer(server, v1API)

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// Start gRPC server
	logger.Log.Info("starting gRPC server...")
	return server.Serve(listen)
}
