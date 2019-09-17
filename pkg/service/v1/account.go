package v1

import (
	"context"
	"fmt"
	"strconv"
	"time"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	"github.com/GameComponent/economy-service/pkg/helper/random"
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

// Authenticate an account
func (s *EconomyServiceServer) Authenticate(ctx context.Context, req *v1.AuthenticateRequest) (*v1.AuthenticateResponse, error) {
	fmt.Println("Authenticate")

	// Check if the user entered to correct credentials
	account, err := s.AccountRepository.GetByEmail(ctx, req.GetEmail())

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	if !checkPasswordHash(req.GetPassword(), account.Hash) {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	// Generate a JWT token
	token, err := s.generateToken(account)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to generate access_token")
	}

	// Generate a save token
	refreshToken, err := random.GenerateRandomString(128)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to generate refresh_token")
	}

	// Calculate the expiration
	refreshExpiration := time.Now().UTC().Add(time.Duration(s.Config.JWTRefreshExpiration) * time.Second)

	// Add the refresh token to the database
	err = s.AccountRepository.CreateRefreshToken(ctx, refreshToken, account.Id, &refreshExpiration)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to generate refresh_token")
	}

	return &v1.AuthenticateResponse{
		AccessToken:  token,
		TokenType:    "Bearer",
		ExpiresIn:    int32(s.Config.JWTExpiration),
		RefreshToken: refreshToken,
	}, nil
}

// Register an account
func (s *EconomyServiceServer) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	fmt.Println("Register")

	// Hash the password
	hash, err := hashPassword(req.GetPassword())
	if err != nil {
		return nil, err
	}

	// Check if user already exists
	exstingAccount, _ := s.AccountRepository.GetByEmail(ctx, req.GetEmail())
	if exstingAccount != nil && exstingAccount.Id != "" {
		return nil, status.Error(codes.AlreadyExists, "user with email already exists")
	}

	// Create the user
	account, err := s.AccountRepository.Create(ctx, req.GetEmail(), hash)
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

// GenerateSecret for an account
func (s *EconomyServiceServer) GenerateSecret(ctx context.Context, req *v1.GenerateSecretRequest) (*v1.GenerateSecretResponse, error) {
	fmt.Println("GenerateSecret")

	if req.GetAccountId() == "" {
		return nil, status.Error(codes.InvalidArgument, "please enter an account_id")
	}

	account, err := s.AccountRepository.Get(ctx, req.GetAccountId())
	if err != nil {
		return nil, status.Error(codes.Internal, "unbable to generate token")
	}

	token, err := s.generateLongLivedToken(account)
	if err != nil {
		return nil, status.Error(codes.Internal, "unbable to generate token")
	}

	return &v1.GenerateSecretResponse{
		Token: token,
	}, nil
}

// Refresh tokens for an account
func (s *EconomyServiceServer) Refresh(ctx context.Context, req *v1.RefreshRequest) (*v1.RefreshResponse, error) {
	fmt.Println("Refresh")

	accountID, err := s.AccountRepository.GetAccountIDFromRefreshToken(ctx, req.Token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unable to find valid refresh_token")
	}

	account, err := s.AccountRepository.Get(ctx, accountID)

	if err != nil {
		return nil, status.Error(codes.Internal, "unable to find account")
	}

	// Generate a JWT token
	token, err := s.generateToken(account)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to generate access_token")
	}

	return &v1.RefreshResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int32(s.Config.JWTExpiration),
	}, nil
}

// GetAccount gets an account
func (s *EconomyServiceServer) GetAccount(ctx context.Context, req *v1.GetAccountRequest) (*v1.GetAccountResponse, error) {
	account, err := s.AccountRepository.Get(ctx, req.GetAccountId())
	if err != nil {
		return nil, err
	}

	// Filter out the account's hash
	if account.Hash != "" {
		account.Hash = ""
	}

	return &v1.GetAccountResponse{
		Account: account,
	}, nil
}

// ListAccount lists accounts
func (s *EconomyServiceServer) ListAccount(ctx context.Context, req *v1.ListAccountRequest) (*v1.ListAccountResponse, error) {
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
	accounts, totalSize, err := s.AccountRepository.List(ctx, limit, offset)
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

// ChangePassword for an account
func (s *EconomyServiceServer) ChangePassword(ctx context.Context, req *v1.ChangePasswordRequest) (*v1.ChangePasswordResponse, error) {
	fmt.Println("ChangePassword")

	// Check if the user entered to correct credentials
	account, err := s.AccountRepository.GetByEmail(ctx, req.GetEmail())

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	if !checkPasswordHash(req.GetPassword(), account.Hash) {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	// Invalidate all existing refresh tokens
	s.AccountRepository.InvalidateRefreshTokens(ctx, account.Id)

	// Hash the password
	hash, err := hashPassword(req.GetNewPassword())
	if err != nil {
		return nil, err
	}

	// Update the account
	updatedAccount, err := s.AccountRepository.Update(ctx, account.Id, hash)
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

// AssignPermission to an account
func (s *EconomyServiceServer) AssignPermission(ctx context.Context, req *v1.AssignPermissionRequest) (*v1.AssignPermissionResponse, error) {
	fmt.Println("AssignPermission")

	if req.GetAccountId() == "" {
		return nil, status.Error(codes.InvalidArgument, "please enter account_id")
	}

	if req.GetPermission() == "" {
		return nil, status.Error(codes.InvalidArgument, "please enter permission")
	}

	account, err := s.AccountRepository.AssignPermission(ctx, req.GetAccountId(), req.GetPermission())
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

// RevokePermission from an account
func (s *EconomyServiceServer) RevokePermission(ctx context.Context, req *v1.RevokePermissionRequest) (*v1.RevokePermissionResponse, error) {
	fmt.Println("RevokePermission")

	account, err := s.AccountRepository.RevokePermission(ctx, req.GetAccountId(), req.GetPermission())
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

func (s *EconomyServiceServer) generateToken(account *v1.Account) (string, error) {
	expirationTime := time.Now().Add(time.Duration(s.Config.JWTExpiration) * time.Second)
	secret := []byte(s.Config.JWTSecret)

	claims := &Claims{
		Subject:     account.Id,
		Email:       account.Email,
		Permissions: account.Permissions,
		StandardClaims: jwt.StandardClaims{
			Audience:  "account",
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

func (s *EconomyServiceServer) generateLongLivedToken(account *v1.Account) (string, error) {
	secret := []byte(s.Config.JWTSecret)

	claims := &Claims{
		Subject:     account.Id,
		Email:       account.Email,
		Permissions: account.Permissions,
		StandardClaims: jwt.StandardClaims{
			Audience: "api",
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
