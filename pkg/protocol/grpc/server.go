package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"google.golang.org/grpc"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	v1service "github.com/GameComponent/economy-service/pkg/service/v1"
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
	// Methods that should always work
	if info.FullMethod == "/v1.EconomyService/Authenticate" {
		return handler(ctx, req)
	}

	if info.FullMethod == "/v1.EconomyService/Register" {
		return handler(ctx, req)
	}

	if info.FullMethod == "/v1.EconomyService/Refresh" {
		return handler(ctx, req)
	}

	// Get the server information
	server, ok := info.Server.(*v1service.EconomyServiceServer)
	if !ok {
		return nil, fmt.Errorf("unable to cast server")
	}

	// Retrieve the secret from the server's config, cast to byte array
	secret := []byte(server.Config.JWTSecret)

	// Check authorization
	_, claims, err := authorize(ctx, secret)

	if err != nil {
		return nil, err
	}

	// Methods that should always work with a valid token
	if info.FullMethod == "/v1.EconomyService/ChangePassword" {
		return handler(ctx, req)
	}

	// Methods that only work if the account has the right permissions
	for _, permission := range claims.Permissions {
		matched, _ := filepath.Match("/v1.EconomyService/"+permission, info.FullMethod)
		if matched {
			return handler(ctx, req)
		}
	}

	return nil, status.Error(codes.PermissionDenied, "Not allowed to execute this method")
}

func authorize(ctx context.Context, secret []byte) (*jwt.Token, *v1service.Claims, error) {
	// Check if metadata is present
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, nil, status.Errorf(codes.InvalidArgument, "Unable to retrieve metadata")
	}

	// Check if authorization header is present
	authorization, ok := md["authorization"]
	if !ok {
		return nil, nil, status.Errorf(codes.Unauthenticated, "Authorization token is not supplied")
	}

	// Extract the key value from the Authorization header
	splits := strings.Split(authorization[0], " ")

	// The token should contain a type and the token
	if len(splits) < 2 {
		return nil, nil, status.Errorf(codes.Unauthenticated, "Invalid token")
	}

	tokenType := strings.ToLower(splits[0])
	tokenString := splits[1]

	// Check if we received a Bearer token
	if tokenType != "bearer" {
		return nil, nil, status.Errorf(codes.Unauthenticated, "Unable to parse this kind of token")
	}

	var claims = &v1service.Claims{}

	// Parse the token and check the encryption method
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	if err != nil {
		return nil, nil, status.Errorf(codes.Unauthenticated, "Unable to parse this token")
	}

	if !token.Valid {
		return nil, nil, status.Errorf(codes.Unauthenticated, "Invalid token")
	}

	// Long lived JWT API tokens should expire and have the audience set to account
	if claims.StandardClaims.Audience == "account" && claims.ExpiresAt != 0 {
		return token, claims, nil
	}

	// Long lived JWT API tokens should not expire and have the audience set to api
	if claims.StandardClaims.Audience == "api" && claims.ExpiresAt == 0 {
		return token, claims, nil
	}

	return nil, nil, status.Errorf(codes.Unauthenticated, "Invalid token")
}
