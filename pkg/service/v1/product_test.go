package v1_test

import (
	"context"
	"fmt"
	"testing"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	mocks "github.com/GameComponent/economy-service/pkg/mocks"
	service "github.com/GameComponent/economy-service/pkg/service/v1"
	assert "github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func TestBuyProductShouldFailIfProductIdIsNil(t *testing.T) {
	config := service.Config{}
	s := service.NewEconomyServiceServer(config)

	req := v1.BuyProductRequest{
		ProductId:          "",
		PriceId:            "",
		PayingStorageId:    "",
		ReceivingStorageId: "",
	}

	result, err := s.BuyProduct(
		context.Background(),
		&req,
	)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "err should not be nil")

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, st.Code(), codes.InvalidArgument, "err status should be codes.InvalidArgument")
}

func TestBuyProductShouldFailIfPriceIdIsNil(t *testing.T) {
	config := service.Config{}
	s := service.NewEconomyServiceServer(config)

	req := v1.BuyProductRequest{
		ProductId:          "product_id",
		PriceId:            "",
		PayingStorageId:    "",
		ReceivingStorageId: "",
	}

	result, err := s.BuyProduct(
		context.Background(),
		&req,
	)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "err should not be nil")

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, st.Code(), codes.InvalidArgument, "err status should be codes.InvalidArgument")
}

func TestBuyProductShouldFailIfPayingStorageIdIsNil(t *testing.T) {
	config := service.Config{}
	s := service.NewEconomyServiceServer(config)

	req := v1.BuyProductRequest{
		ProductId:          "product_id",
		PriceId:            "price_id",
		PayingStorageId:    "",
		ReceivingStorageId: "",
	}

	result, err := s.BuyProduct(
		context.Background(),
		&req,
	)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "err should not be nil")

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, st.Code(), codes.InvalidArgument, "err status should be codes.InvalidArgument")
}

func TestBuyProductShouldFailIfReceivingStorageIdIsNil(t *testing.T) {
	config := service.Config{}
	s := service.NewEconomyServiceServer(config)

	req := v1.BuyProductRequest{
		ProductId:          "product_id",
		PriceId:            "price_id",
		PayingStorageId:    "paying_storage_id",
		ReceivingStorageId: "",
	}

	result, err := s.BuyProduct(
		context.Background(),
		&req,
	)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "err should not be nil")

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, st.Code(), codes.InvalidArgument, "err status should be codes.InvalidArgument")
}

func TestBuyProductShouldFailIfProductIsNotFound(t *testing.T) {
	mockProductRepository := mocks.ProductRepository{}

	// Get returns nil and an error
	mockProductRepository.On("Get", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("wat"))

	// Create the service and inject the mocked ProductRepository
	config := service.Config{
		ProductRepository: &mockProductRepository,
	}
	s := service.NewEconomyServiceServer(config)

	req := v1.BuyProductRequest{
		ProductId:          "product_id",
		PriceId:            "price_id",
		PayingStorageId:    "paying_storage_id",
		ReceivingStorageId: "receiving_storage_id",
	}

	result, err := s.BuyProduct(
		context.Background(),
		&req,
	)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "err should not be nil")

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, st.Code(), codes.NotFound, "err status should be codes.NotFound")
}

func TestBuyProductShouldFailIfPriceIsNotPartOfProduct(t *testing.T) {
	mockProductRepository := mocks.ProductRepository{}
	mockProductRepository.On("Get", mock.Anything, mock.Anything).Return(&v1.Product{}, nil)

	// Create the service and inject the mocked ProductRepository
	config := service.Config{
		ProductRepository: &mockProductRepository,
	}
	s := service.NewEconomyServiceServer(config)

	req := v1.BuyProductRequest{
		ProductId:          "product_id",
		PriceId:            "price_id",
		PayingStorageId:    "paying_storage_id",
		ReceivingStorageId: "receiving_storage_id",
	}

	result, err := s.BuyProduct(
		context.Background(),
		&req,
	)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "err should not be nil")

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, st.Code(), codes.NotFound, "err status should be codes.NotFound")
}

func TestBuyProductShouldFailIfUnableToRetrievePayingStorage(t *testing.T) {
	// Mock the ProductRepository
	mockProductRepository := mocks.ProductRepository{}
	mockPrice := v1.Price{
		Id: "price_id",
	}
	mockProduct := v1.Product{
		Id: "product_id",
		Prices: []*v1.Price{
			&mockPrice,
		},
	}
	mockProductRepository.On("Get", mock.Anything, mock.Anything).Return(&mockProduct, nil)

	// Mock the StorageRepository
	mockStorageRepository := mocks.StorageRepository{}
	mockStorageRepository.On("Get", mock.Anything, mock.Anything).Return(nil, nil)

	// Create the service and inject the mocked ProductRepository
	config := service.Config{
		ProductRepository: &mockProductRepository,
		StorageRepository: &mockStorageRepository,
	}
	s := service.NewEconomyServiceServer(config)

	req := v1.BuyProductRequest{
		ProductId:          "product_id",
		PriceId:            "price_id",
		PayingStorageId:    "paying_storage_id",
		ReceivingStorageId: "receiving_storage_id",
	}

	result, err := s.BuyProduct(
		context.Background(),
		&req,
	)

	assert.Nil(t, result, "result should be nil")
	assert.NotNil(t, err, "err should not be nil")

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, st.Code(), codes.NotFound, "err status should be codes.NotFound")
}
