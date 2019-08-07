package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"

	"google.golang.org/grpc"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	jwt "github.com/dgrijalva/jwt-go"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// RunServer runs gRPC service to publish Economy service
func RunServer(ctx context.Context, v1API v1.EconomyServiceServer, logger *zap.Logger, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// Register service
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryInterceptor)),
	)

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
	logger.Info("starting gRPC server", zap.String("port", port))
	return server.Serve(listen)
}

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod == "/v1.EconomyService/Authenticate" {
		return handler(ctx, req)
	}

	if info.FullMethod == "/v1.EconomyService/Register" {
		return handler(ctx, req)
	}

	// Check authorization
	if err := authorize(ctx); err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

func authorize(ctx context.Context) error {
	// Check if metadata is present
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.InvalidArgument, "Unable to retrieve metadata")
	}

	// Check if authorization header is present
	authorization, ok := md["authorization"]
	if !ok {
		return status.Errorf(codes.Unauthenticated, "Authorization token is not supplied")
	}

	// Extract the key value from the Authorization header
	splits := strings.Split(authorization[0], " ")
	tokenType := strings.ToLower(splits[0])
	tokenString := splits[1]

	// Check if we received a Bearer token
	if tokenType != "bearer" {
		return status.Errorf(codes.Unauthenticated, "Unable to parse this kind of token")
	}

	// TODO: get secret another way
	var secret = []byte("my_secret_key")

	// Parse the token and check the encryption method
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "Unable to parse this token")
	}

	if !token.Valid {
		return status.Errorf(codes.Unauthenticated, "Invalid token")
	}

	return nil
}
