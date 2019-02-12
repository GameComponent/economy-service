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

	// Add item to the databased return the generated UUID
	lastInsertUuid := ""
	err := s.db.QueryRowContext(
		ctx,
		`INSERT INTO currency(name) VALUES ($1) RETURNING id`,
		req.GetName(),
	).Scan(&lastInsertUuid)

	if err != nil {
		return nil, err
	}

	// Generate the object based on the generated id and the requested name
	currency := &v1.Currency{
		Id:   lastInsertUuid,
		Name: req.GetName(),
	}

	return &v1.CreateCurrencyResponse{
		Api:      apiVersion,
		Currency: currency,
	}, nil
}
