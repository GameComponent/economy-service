package v1

import (
	"context"
	"fmt"
	"strconv"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// ListCurrency lists currencies
func (s *EconomyServiceServer) ListCurrency(ctx context.Context, req *v1.ListCurrencyRequest) (*v1.ListCurrencyResponse, error) {
	fmt.Println("ListCurrency")

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

	// Get the currencies from the repository
	currencies, totalSize, err := s.CurrencyRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to retrieve currency list")
	}

	// Determine if there is a next page
	var nextPageToken string
	if totalSize > (offset + limit) {
		nextPage := int32(parsedToken) + 1
		nextPageToken = strconv.Itoa(int(nextPage))
	}

	return &v1.ListCurrencyResponse{
		Currencies:    currencies,
		TotalSize:     totalSize,
		NextPageToken: nextPageToken,
	}, nil
}

// CreateCurrency creates a new currency
func (s *EconomyServiceServer) CreateCurrency(ctx context.Context, req *v1.CreateCurrencyRequest) (*v1.CreateCurrencyResponse, error) {
	fmt.Println("CreateCurrency")

	// Check the name
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "no name given")
	}

	// Check the short_name
	if req.GetShortName() == "" {
		return nil, status.Error(codes.InvalidArgument, "no short_name given")
	}

	// Check the symbol
	if req.GetSymbol() == "" {
		return nil, status.Error(codes.InvalidArgument, "no symbol given")
	}

	currency, err := s.CurrencyRepository.Create(
		ctx,
		req.GetName(),
		req.GetShortName(),
		req.GetSymbol(),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to create currency")
	}

	return &v1.CreateCurrencyResponse{
		Currency: currency,
	}, nil
}

// UpdateCurrency updates an existing currency
func (s *EconomyServiceServer) UpdateCurrency(ctx context.Context, req *v1.UpdateCurrencyRequest) (*v1.UpdateCurrencyResponse, error) {
	fmt.Println("UpdateCurrency")

	// Check the id
	if req.GetCurrencyId() == "" {
		return nil, status.Error(codes.InvalidArgument, "no currency_id given")
	}

	currency, err := s.CurrencyRepository.Update(
		ctx,
		req.GetCurrencyId(),
		req.GetName(),
		req.GetShortName(),
		req.GetSymbol(),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to update currency")
	}

	return &v1.UpdateCurrencyResponse{
		Currency: currency,
	}, nil
}

// GetCurrency gets a currency
func (s *EconomyServiceServer) GetCurrency(ctx context.Context, req *v1.GetCurrencyRequest) (*v1.GetCurrencyResponse, error) {
	fmt.Println("GetCurrency")

	currency, err := s.CurrencyRepository.Get(ctx, req.GetCurrencyId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "currency not found")
	}

	return &v1.GetCurrencyResponse{
		Currency: currency,
	}, nil
}
