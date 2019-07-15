package v1

import (
	"context"
	"fmt"
	"strconv"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *economyServiceServer) ListCurrency(ctx context.Context, req *v1.ListCurrencyRequest) (*v1.ListCurrencyResponse, error) {
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
	currencies, totalSize, err := s.currencyRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, err
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

func (s *economyServiceServer) CreateCurrency(ctx context.Context, req *v1.CreateCurrencyRequest) (*v1.CreateCurrencyResponse, error) {
	fmt.Println("CreateCurrency")

	// Check the name
	if len(req.GetName()) == 0 {
		return nil, fmt.Errorf("currency should have a name")
	}

	// Check the short_name
	if len(req.GetShortName()) == 0 {
		return nil, fmt.Errorf("currency should have a short_name")
	}

	// Check the symbol
	if len(req.GetSymbol()) == 0 {
		return nil, fmt.Errorf("currency should have a symbol")
	}

	currency, err := s.currencyRepository.Create(
		ctx,
		req.GetName(),
		req.GetShortName(),
		req.GetSymbol(),
	)
	if err != nil {
		return nil, err
	}

	return &v1.CreateCurrencyResponse{
		Currency: currency,
	}, nil
}

func (s *economyServiceServer) UpdateCurrency(ctx context.Context, req *v1.UpdateCurrencyRequest) (*v1.UpdateCurrencyResponse, error) {
	fmt.Println("UpdateCurrency")

	// Check the name
	if len(req.GetName()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no name given")
	}

	// Check the short_name
	if len(req.GetShortName()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no short_name given")
	}

	// Check the symbol
	if len(req.GetSymbol()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "no symbol given")
	}

	currency, err := s.currencyRepository.Update(
		ctx,
		req.GetCurrencyId(),
		req.GetName(),
		req.GetShortName(),
		req.GetSymbol(),
	)
	if err != nil {
		return nil, err
	}

	return &v1.UpdateCurrencyResponse{
		Currency: currency,
	}, nil
}

func (s *economyServiceServer) GetCurrency(ctx context.Context, req *v1.GetCurrencyRequest) (*v1.GetCurrencyResponse, error) {
	fmt.Println("GetCurrency")

	currency, err := s.currencyRepository.Get(ctx, req.GetCurrencyId())
	if err != nil {
		return nil, err
	}

	return &v1.GetCurrencyResponse{
		Currency: currency,
	}, nil
}
