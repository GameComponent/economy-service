package grpc

import (
  "context"
  "log"
  "net"
  "os"
  "os/signal"

  "google.golang.org/grpc"

  "github.com/GameComponent/economy-service/pkg/api/v1"
)

// RunServer runs gRPC service to publish Economy service
func RunServer(ctx context.Context, v1API v1.EconomyServiceServer, port string) error {
  listen, err := net.Listen("tcp", ":"+port)
  if err != nil {
    return err
  }

  // register service
  server := grpc.NewServer()
  v1.RegisterEconomyServiceServer(server, v1API)

  // graceful shutdown
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

  // start gRPC server
  log.Println("starting gRPC server...")
  return server.Serve(listen)
}
