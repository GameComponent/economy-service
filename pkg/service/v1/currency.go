package v1

import (
	"context"
	"fmt"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
)

func (s *economyServiceServer) ListCurrency(ctx context.Context, req *v1.ListCurrencyRequest) (*v1.ListCurrencyResponse, error) {
	fmt.Println("ListCurrency")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	return &v1.ListCurrencyResponse{
		Api:        apiVersion,
		Currencies: []*v1.Currency{},
	}, nil
}

func (s *economyServiceServer) CreateCurrency(ctx context.Context, req *v1.CreateCurrencyRequest) (*v1.CreateCurrencyResponse, error) {
	fmt.Println("CreateCurrency")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	currency, err := s.currencyRepository.Create(ctx, req.GetName())
	if err != nil {
		return nil, err
	}

	return &v1.CreateCurrencyResponse{
		Api:      apiVersion,
		Currency: currency,
	}, nil
}

func (s *economyServiceServer) GetCurrency(ctx context.Context, req *v1.GetCurrencyRequest) (*v1.GetCurrencyResponse, error) {
	fmt.Println("GetCurrency")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	currency, err := s.currencyRepository.Get(ctx, req.GetCurrencyId())
	if err != nil {
		return nil, err
	}

	return &v1.GetCurrencyResponse{
		Api:      apiVersion,
		Currency: currency,
	}, nil
}
