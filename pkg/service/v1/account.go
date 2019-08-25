package v1

import (
	"context"
	"fmt"
	"strconv"
	"time"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	jwt "github.com/dgrijalva/jwt-go"
	bcrypt "golang.org/x/crypto/bcrypt"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Claims for the JWT token
type Claims struct {
	Subject     string   `json:"sub"`
	Email       string   `json:"email"`
	Permissions []string `json:"permissions"`
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

	// Check if the user entered to correct credentials
	account, err := s.accountRepository.GetByEmail(ctx, req.GetEmail())

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	if !checkPasswordHash(req.GetPassword(), account.Hash) {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	// Generate a JWT token
	token, err := s.generateToken(account)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to generate token")
	}

	return &v1.AuthenticateResponse{
		Token: token,
	}, nil
}

func (s *economyServiceServer) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	fmt.Println("Register")

	// Hash the password
	hash, err := hashPassword(req.GetPassword())
	if err != nil {
		return nil, err
	}

	// Check if user already exists
	exstingAccount, _ := s.accountRepository.GetByEmail(ctx, req.GetEmail())
	if exstingAccount != nil {
		return nil, status.Error(codes.AlreadyExists, "user with email already exists")
	}

	// Create the user
	account, err := s.accountRepository.Create(ctx, req.GetEmail(), hash)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to create account")
	}

	// Generate a JWT token
	token, err := s.generateToken(account)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to generate token")
	}

	return &v1.RegisterResponse{
		Token: token,
	}, nil
}

func (s *economyServiceServer) ListAccount(ctx context.Context, req *v1.ListAccountRequest) (*v1.ListAccountResponse, error) {
	fmt.Println("ListAccount")

	// Parse the page token
	var parsedToken int64
	parsedToken, _ = strconv.ParseInt(req.GetPageToken(), 10, 32)

	// Get the limit
	limit := req.GetPageSize()
	if limit == 0 {
		limit = 100
	}

	// Get the offset
	offset := int32(0)
	if len(req.GetPageToken()) > 0 {
		offset = int32(parsedToken) * limit
	}

	// Get the accounts from the repository
	accounts, totalSize, err := s.accountRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to retrieve account list")
	}

	// Determine if there is a next page
	var nextPageToken string
	if totalSize > (offset + limit) {
		nextPage := int32(parsedToken) + 1
		nextPageToken = strconv.Itoa(int(nextPage))
	}

	return &v1.ListAccountResponse{
		Accounts:      accounts,
		TotalSize:     totalSize,
		NextPageToken: nextPageToken,
	}, nil
}

func (s *economyServiceServer) ChangePassword(ctx context.Context, req *v1.ChangePasswordRequest) (*v1.ChangePasswordResponse, error) {
	fmt.Println("ChangePassword")

	// Check if the user entered to correct credentials
	account, err := s.accountRepository.GetByEmail(ctx, req.GetEmail())

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	if !checkPasswordHash(req.GetPassword(), account.Hash) {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	// Hash the password
	hash, err := hashPassword(req.GetNewPassword())
	if err != nil {
		return nil, err
	}

	// Update the account
	updatedAccount, err := s.accountRepository.Update(ctx, account.Id, hash)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to update account")
	}

	// Generate a JWT token
	token, err := s.generateToken(updatedAccount)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to generate token")
	}

	return &v1.ChangePasswordResponse{
		Token: token,
	}, nil
}

func (s *economyServiceServer) AssignPermission(ctx context.Context, req *v1.AssignPermissionRequest) (*v1.AssignPermissionResponse, error) {
	fmt.Println("AssignPermission")

	account, err := s.accountRepository.AssignPermission(ctx, req.GetAccountId(), req.GetPermission())
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to assign permission")
	}

	// Filter out the account's hash
	if account.Hash != "" {
		account.Hash = ""
	}

	return &v1.AssignPermissionResponse{
		Account: account,
	}, nil
}

func (s *economyServiceServer) RevokePermission(ctx context.Context, req *v1.RevokePermissionRequest) (*v1.RevokePermissionResponse, error) {
	fmt.Println("RevokePermission")

	account, err := s.accountRepository.RevokePermission(ctx, req.GetAccountId(), req.GetPermission())
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to revoke permission")
	}

	// Filter out the account's hash
	if account.Hash != "" {
		account.Hash = ""
	}

	return &v1.RevokePermissionResponse{
		Account: account,
	}, nil
}

func (s *economyServiceServer) generateToken(account *v1.Account) (string, error) {
	expirationTime := time.Now().Add(time.Duration(s.config.JWTExpiration) * time.Second)
	secret := []byte(s.config.JWTSecret)

	claims := &Claims{
		Subject:     account.Id,
		Email:       account.Email,
		Permissions: account.Permissions,
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
		return "", fmt.Errorf("unable to sign token")
	}

	return tokenString, nil
}
