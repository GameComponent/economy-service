package v1

import (
	"context"
	"fmt"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// GetPrice get a price
func (s *EconomyServiceServer) GetPrice(ctx context.Context, req *v1.GetPriceRequest) (*v1.GetPriceResponse, error) {
	fmt.Println("GetPrice")

	price, err := s.PriceRepository.Get(ctx, req.GetPriceId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "price not found")
	}

	return &v1.GetPriceResponse{
		Price: price,
	}, nil
}

// CreatePrice creates a new price
func (s *EconomyServiceServer) CreatePrice(ctx context.Context, req *v1.CreatePriceRequest) (*v1.CreatePriceResponse, error) {
	fmt.Println("CreatePrice")

	if req.GetProductId() == "" {
		return nil, status.Error(codes.InvalidArgument, "no product_id given")
	}

	// Add the price to the databased return the generated UUID
	price, err := s.PriceRepository.Create(ctx, req.GetProductId())
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to create price")
	}

	return &v1.CreatePriceResponse{
		Price: price,
	}, nil
}

// AttachPriceCurrency attaches a currency to a price
func (s *EconomyServiceServer) AttachPriceCurrency(ctx context.Context, req *v1.AttachPriceCurrencyRequest) (*v1.AttachPriceCurrencyResponse, error) {
	fmt.Println("AttachPriceCurrency")

	price, err := s.PriceRepository.AttachPriceCurrency(
		ctx,
		req.GetPriceId(),
		req.GetCurrencyId(),
		req.GetAmount(),
	)

	if err != nil {
		return nil, status.Error(codes.Internal, "unable to attach currency to price")
	}

	return &v1.AttachPriceCurrencyResponse{
		Price: price,
	}, nil
}

// DetachPriceCurrency detaches a currency from a price
func (s *EconomyServiceServer) DetachPriceCurrency(ctx context.Context, req *v1.DetachPriceCurrencyRequest) (*v1.DetachPriceCurrencyResponse, error) {
	fmt.Println("DetachPriceCurrency")

	price, err := s.PriceRepository.DetachPriceCurrency(
		ctx,
		req.GetPriceCurrencyId(),
	)

	if err != nil {
		return nil, status.Error(codes.Internal, "unable to detach currency from price")
	}

	return &v1.DetachPriceCurrencyResponse{
		Price: price,
	}, nil
}

// AttachPriceItem attaches an item to a price
func (s *EconomyServiceServer) AttachPriceItem(ctx context.Context, req *v1.AttachPriceItemRequest) (*v1.AttachPriceItemResponse, error) {
	fmt.Println("AttachPriceItem")

	price, err := s.PriceRepository.AttachPriceItem(
		ctx,
		req.GetPriceId(),
		req.GetItemId(),
		req.GetAmount(),
	)

	if err != nil {
		return nil, status.Error(codes.Internal, "unable to attach item to price")
	}

	return &v1.AttachPriceItemResponse{
		Price: price,
	}, nil
}

// DetachPriceItem detaches an item from a price
func (s *EconomyServiceServer) DetachPriceItem(ctx context.Context, req *v1.DetachPriceItemRequest) (*v1.DetachPriceItemResponse, error) {
	fmt.Println("DetachPriceItem")

	price, err := s.PriceRepository.DetachPriceCurrency(
		ctx,
		req.GetPriceItemId(),
	)

	if err != nil {
		return nil, status.Error(codes.Internal, "unable to detach item from price")
	}

	return &v1.DetachPriceItemResponse{
		Price: price,
	}, nil
}

// DeletePrice deletes a price
func (s *EconomyServiceServer) DeletePrice(ctx context.Context, req *v1.DeletePriceRequest) (*v1.DeletePriceResponse, error) {
	fmt.Println("DeletePrice")

	success, err := s.PriceRepository.Delete(
		ctx,
		req.GetPriceId(),
	)

	if err != nil {
		return nil, status.Error(codes.NotFound, "price not found")
	}

	return &v1.DeletePriceResponse{
		Success: success,
	}, nil
}
