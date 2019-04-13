package v1

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var secret = []byte("my_secret_key")

// Claims for the JWT token
type Claims struct {
	Subject string `json:"sub"`
	Email   string `json:"email"`
	jwt.StandardClaims
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *economyServiceServer) Authenticate(ctx context.Context, req *v1.AuthenticateRequest) (*v1.AuthenticateResponse, error) {
	fmt.Println("Authenticate")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// Check if the user entered to correct credentials
	account := s.accountRepository.Get(ctx, req.GetEmail())

	if account == nil {
		return nil, fmt.Errorf("Invalid account credentials")
	}

	if !checkPasswordHash(req.GetPassword(), account.Hash) {
		return nil, fmt.Errorf("Invalid account credentials2")
	}

	expirationTime := time.Now().Add(2 * time.Hour)
	claims := &Claims{
		Subject: account.ID,
		Email:   account.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return nil, fmt.Errorf("Unable to create token")
	}

	return &v1.AuthenticateResponse{
		Api:   apiVersion,
		Token: tokenString,
	}, nil
}

func (s *economyServiceServer) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	fmt.Println("Register")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	// Hash the password
	hash, err := hashPassword(req.GetPassword())
	if err != nil {
		return nil, err
	}

	// Check if the user entered to correct credentials
	account := s.accountRepository.Create(ctx, req.GetEmail(), hash)
	if account == nil {
		return nil, fmt.Errorf("Unable to create account")
	}

	expirationTime := time.Now().Add(2 * time.Hour)
	claims := &Claims{
		Subject: account.ID,
		Email:   account.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return nil, fmt.Errorf("Unable to create token")
	}

	return &v1.RegisterResponse{
		Api:   apiVersion,
		Token: tokenString,
	}, nil
}
