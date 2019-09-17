package v1

import (
	"context"
	"fmt"
	"strconv"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// CreateProduct creates a new product
func (s *EconomyServiceServer) CreateProduct(ctx context.Context, req *v1.CreateProductRequest) (*v1.CreateProductResponse, error) {
	fmt.Println("CreateProduct")

	// Add product to the databased return the generated UUID
	product, err := s.ProductRepository.Create(ctx, req.GetName())
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to create product")
	}

	return &v1.CreateProductResponse{
		Product: product,
	}, nil
}

// UpdateProduct updates a product
func (s *EconomyServiceServer) UpdateProduct(ctx context.Context, req *v1.UpdateProductRequest) (*v1.UpdateProductResponse, error) {
	fmt.Println("UpdateProduct")

	if req.GetProductId() == "" {
		return nil, status.Error(codes.InvalidArgument, "no product_id given")
	}

	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "no name given")
	}

	// Add product to the databased return the generated UUID
	product, err := s.ProductRepository.Update(ctx, req.GetProductId(), req.GetName())
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to update product")
	}

	return &v1.UpdateProductResponse{
		Product: product,
	}, nil
}

// ListProduct lists products
func (s *EconomyServiceServer) ListProduct(ctx context.Context, req *v1.ListProductRequest) (*v1.ListProductResponse, error) {
	fmt.Println("ListProduct")

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

	// Get the products
	products, totalSize, err := s.ProductRepository.List(ctx, limit, offset)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to retrieve product list")
	}

	// Determine if there is a next page
	var nextPageToken string
	if totalSize > (offset + limit) {
		nextPage := int32(parsedToken) + 1
		nextPageToken = strconv.Itoa(int(nextPage))
	}

	return &v1.ListProductResponse{
		Products:      products,
		TotalSize:     totalSize,
		NextPageToken: nextPageToken,
	}, nil
}

// GetProduct gets a product
func (s *EconomyServiceServer) GetProduct(ctx context.Context, req *v1.GetProductRequest) (*v1.GetProductResponse, error) {
	product, err := s.ProductRepository.Get(ctx, req.GetProductId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "product not found")
	}

	return &v1.GetProductResponse{
		Product: product,
	}, nil
}

// AttachItem attaches an item to a product
func (s *EconomyServiceServer) AttachItem(ctx context.Context, req *v1.AttachItemRequest) (*v1.AttachItemResponse, error) {
	fmt.Println("AttachItem")

	// Add product to the databased return the generated UUID
	product, err := s.ProductRepository.AttachItem(
		ctx,
		req.GetProductId(),
		req.GetItemId(),
		req.GetAmount(),
	)

	if err != nil {
		return nil, status.Error(codes.Internal, "unable to attach item to product")
	}

	return &v1.AttachItemResponse{
		Product: product,
	}, nil
}

// DetachItem detaches an item from a product
func (s *EconomyServiceServer) DetachItem(ctx context.Context, req *v1.DetachItemRequest) (*v1.DetachItemResponse, error) {
	fmt.Println("DetachItem")

	// Add product to the databased return the generated UUID
	product, err := s.ProductRepository.DetachItem(
		ctx,
		req.GetProductItemId(),
	)

	if err != nil {
		return nil, status.Error(codes.Internal, "unable to detach item from product")
	}

	return &v1.DetachItemResponse{
		Product: product,
	}, nil
}

// AttachCurrency attaches a currency to a product
func (s *EconomyServiceServer) AttachCurrency(ctx context.Context, req *v1.AttachCurrencyRequest) (*v1.AttachCurrencyResponse, error) {
	fmt.Println("AttachCurrency")

	// Add product to the databased return the generated UUID
	product, err := s.ProductRepository.AttachCurrency(
		ctx,
		req.GetProductId(),
		req.GetCurrencyId(),
		req.GetAmount(),
	)

	if err != nil {
		return nil, status.Error(codes.Internal, "unable to attach currency to product")
	}

	return &v1.AttachCurrencyResponse{
		Product: product,
	}, nil
}

// DetachCurrency detaches a currency from a product
func (s *EconomyServiceServer) DetachCurrency(ctx context.Context, req *v1.DetachCurrencyRequest) (*v1.DetachCurrencyResponse, error) {
	fmt.Println("DetachCurrency")

	// Add product to the databased return the generated UUID
	product, err := s.ProductRepository.DetachCurrency(
		ctx,
		req.GetProductCurrencyId(),
	)

	if err != nil {
		return nil, status.Error(codes.Internal, "unable to detach item from product")
	}

	return &v1.DetachCurrencyResponse{
		Product: product,
	}, nil
}

// ListProductPrice lists all prices for the product
func (s *EconomyServiceServer) ListProductPrice(ctx context.Context, req *v1.ListProductPriceRequest) (*v1.ListProductPriceResponse, error) {
	fmt.Println("ListProductPrice")

	if req.GetProductId() == "" {
		return nil, status.Error(codes.InvalidArgument, "no product_id given")
	}

	prices, err := s.ProductRepository.ListPrice(ctx, req.GetProductId())
	if err != nil {
		return nil, err
	}

	return &v1.ListProductPriceResponse{
		Prices: prices,
	}, nil
}

// BuyProduct buys the product
func (s *EconomyServiceServer) BuyProduct(ctx context.Context, req *v1.BuyProductRequest) (*v1.BuyProductResponse, error) {
	fmt.Println("BuyProduct")

	if req.GetProductId() == "" {
		return nil, status.Error(codes.InvalidArgument, "no product_id given")
	}

	if req.GetPriceId() == "" {
		return nil, status.Error(codes.InvalidArgument, "no price_id given")
	}

	if req.GetReceivingStorageId() == "" {
		return nil, status.Error(codes.InvalidArgument, "no receiving_storage_id given")
	}

	if req.GetPayingStorageId() == "" {
		return nil, status.Error(codes.InvalidArgument, "no paying_storage_id given")
	}

	// Get the Product
	product, err := s.ProductRepository.Get(ctx, req.GetProductId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "product not found")
	}

	// Turn the Price slice into a map
	productPrices := product.Prices
	productPricesMap := map[string]*v1.Price{}
	for _, priceItem := range productPrices {
		productPricesMap[priceItem.Id] = priceItem
	}

	// Check if the price is part of the Product
	price := productPricesMap[req.GetPriceId()]
	if price == nil || price.Id == "" {
		return nil, status.Error(codes.NotFound, "price not found in product")
	}

	// Get the paying Storage
	payingStorage, err := s.StorageRepository.Get(ctx, req.GetPayingStorageId())
	if payingStorage == nil || payingStorage.Id == "" {
		return nil, status.Error(codes.NotFound, "paying_storage_id not found")
	}

	// Get the receiving Storage
	receivingStorage := &v1.Storage{}

	// Check if the paying and receiving Storage are equal
	if req.GetPayingStorageId() == req.GetReceivingStorageId() {
		receivingStorage = payingStorage
	}

	// Receiving Storage is different lets get it
	if req.GetPayingStorageId() != req.GetReceivingStorageId() {
		receivingStorage, err = s.StorageRepository.Get(ctx, req.GetReceivingStorageId())
		if err != nil {
			return nil, err
		}
	}

	if receivingStorage.Id == "" {
		return nil, status.Error(codes.NotFound, "receiving_storage_id not found")
	}

	// Determine if there is enough of the Currency in the paying Storage
	for _, priceCurrency := range price.Currencies {
		hasEnoughOfCurrency := false

		for _, storageCurrency := range payingStorage.Currencies {
			if storageCurrency.Currency.Id != priceCurrency.Currency.Id {
				continue
			}

			if storageCurrency.Amount >= priceCurrency.Amount {
				hasEnoughOfCurrency = true
			}
		}

		if !hasEnoughOfCurrency {
			return nil, status.Error(
				codes.Aborted,
				fmt.Sprintf("not enough of currency %s in the storage", priceCurrency.Currency.Id),
			)
		}
	}

	// Determine if there are enough Items in the paying Storage
	for _, priceItem := range price.Items {
		remainingItems := priceItem.Amount

		for _, storageItem := range payingStorage.Items {
			if storageItem.Item.Id != priceItem.Item.Id {
				continue
			}

			// Calculate the amount of items
			// Non-stackable items should count as 1
			amount := storageItem.Amount
			if amount == 0 && !storageItem.Item.Stackable {
				amount = 1
			}

			remainingItems = remainingItems - amount
		}

		if remainingItems > 0 {
			return nil, status.Error(
				codes.Aborted,
				fmt.Sprintf("not enough of items %s in the storage", priceItem.Item.Id),
			)
		}
	}

	_, err = s.ProductRepository.BuyProduct(
		ctx,
		product,
		price,
		receivingStorage,
		payingStorage,
	)
	if err != nil {
		return nil, err
	}

	return &v1.BuyProductResponse{
		Product: product,
	}, nil
}
