package v1

import (
	"context"
	"fmt"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
)

func (s *economyServiceServer) GetPrice(ctx context.Context, req *v1.GetPriceRequest) (*v1.GetPriceResponse, error) {
	fmt.Println("GetPrice")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	price, err := s.priceRepository.Get(ctx, req.GetPriceId())
	if err != nil {
		return nil, err
		// return nil, fmt.Errorf("unable to retrieve price")
	}

	return &v1.GetPriceResponse{
		Api:   apiVersion,
		Price: price,
	}, nil
}

func (s *economyServiceServer) CreatePrice(ctx context.Context, req *v1.CreatePriceRequest) (*v1.CreatePriceResponse, error) {
	fmt.Println("CreatePrice")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	if req.GetProductId() == "" {
		return nil, fmt.Errorf("please speficy a product id")
	}

	// Add the price to the databased return the generated UUID
	price, err := s.priceRepository.Create(ctx, req.GetProductId())
	if err != nil {
		return nil, err
	}

	return &v1.CreatePriceResponse{
		Api:   apiVersion,
		Price: price,
	}, nil
}

func (s *economyServiceServer) AttachPriceCurrency(ctx context.Context, req *v1.AttachPriceCurrencyRequest) (*v1.AttachPriceCurrencyResponse, error) {
	fmt.Println("AttachPriceCurrency")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	price, err := s.priceRepository.AttachPriceCurrency(
		ctx,
		req.GetPriceId(),
		req.GetCurrencyId(),
		req.GetAmount(),
	)

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve price")
	}

	return &v1.AttachPriceCurrencyResponse{
		Api:   apiVersion,
		Price: price,
	}, nil
}

func (s *economyServiceServer) DetachPriceCurrency(ctx context.Context, req *v1.DetachPriceCurrencyRequest) (*v1.DetachPriceCurrencyResponse, error) {
	fmt.Println("DetachPriceCurrency")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	price, err := s.priceRepository.DetachPriceCurrency(
		ctx,
		req.GetPriceCurrencyId(),
	)

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve price")
	}

	return &v1.DetachPriceCurrencyResponse{
		Api:   apiVersion,
		Price: price,
	}, nil
}

func (s *economyServiceServer) AttachPriceItem(ctx context.Context, req *v1.AttachPriceItemRequest) (*v1.AttachPriceItemResponse, error) {
	fmt.Println("AttachPriceItem")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	price, err := s.priceRepository.AttachPriceItem(
		ctx,
		req.GetPriceId(),
		req.GetItemId(),
		req.GetAmount(),
	)

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve price")
	}

	return &v1.AttachPriceItemResponse{
		Api:   apiVersion,
		Price: price,
	}, nil
}

func (s *economyServiceServer) DetachPriceItem(ctx context.Context, req *v1.DetachPriceItemRequest) (*v1.DetachPriceItemResponse, error) {
	fmt.Println("DetachPriceItem")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	price, err := s.priceRepository.DetachPriceCurrency(
		ctx,
		req.GetPriceItemId(),
	)

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve price")
	}

	return &v1.DetachPriceItemResponse{
		Api:   apiVersion,
		Price: price,
	}, nil
}

func (s *economyServiceServer) DeletePrice(ctx context.Context, req *v1.DeletePriceRequest) (*v1.DeletePriceResponse, error) {
	fmt.Println("DeletePrice")

	// check if the API version requested by client is supported by server
	if err := s.checkAPI(req.Api); err != nil {
		return nil, err
	}

	success, err := s.priceRepository.Delete(
		ctx,
		req.GetPriceId(),
	)

	if err != nil {
		return nil, fmt.Errorf("unable to delete price")
	}

	return &v1.DeletePriceResponse{
		Api:     apiVersion,
		Success: success,
	}, nil
}
