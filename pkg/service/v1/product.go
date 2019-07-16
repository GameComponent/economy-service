package v1

import (
	"context"
	"fmt"
	"strconv"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func (s *economyServiceServer) CreateProduct(ctx context.Context, req *v1.CreateProductRequest) (*v1.CreateProductResponse, error) {
	fmt.Println("CreateProduct")

	// Add product to the databased return the generated UUID
	product, err := s.productRepository.Create(ctx, req.GetName())
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to create product")
	}

	return &v1.CreateProductResponse{
		Product: product,
	}, nil
}

func (s *economyServiceServer) UpdateProduct(ctx context.Context, req *v1.UpdateProductRequest) (*v1.UpdateProductResponse, error) {
	fmt.Println("UpdateProduct")

	if req.GetProductId() == "" {
		return nil, status.Error(codes.InvalidArgument, "no product_id given")
	}

	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "no name given")
	}

	// Add product to the databased return the generated UUID
	product, err := s.productRepository.Update(ctx, req.GetProductId(), req.GetName())
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to update product")
	}

	return &v1.UpdateProductResponse{
		Product: product,
	}, nil
}

func (s *economyServiceServer) ListProduct(ctx context.Context, req *v1.ListProductRequest) (*v1.ListProductResponse, error) {
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
	products, totalSize, err := s.productRepository.List(ctx, limit, offset)
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

func (s *economyServiceServer) GetProduct(ctx context.Context, req *v1.GetProductRequest) (*v1.GetProductResponse, error) {
	product, err := s.productRepository.Get(ctx, req.GetProductId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "product not found")
	}

	return &v1.GetProductResponse{
		Product: product,
	}, nil
}

func (s *economyServiceServer) AttachItem(ctx context.Context, req *v1.AttachItemRequest) (*v1.AttachItemResponse, error) {
	fmt.Println("AttachItem")

	// Add product to the databased return the generated UUID
	product, err := s.productRepository.AttachItem(
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

func (s *economyServiceServer) DetachItem(ctx context.Context, req *v1.DetachItemRequest) (*v1.DetachItemResponse, error) {
	fmt.Println("DetachItem")

	// Add product to the databased return the generated UUID
	product, err := s.productRepository.DetachItem(
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

func (s *economyServiceServer) AttachCurrency(ctx context.Context, req *v1.AttachCurrencyRequest) (*v1.AttachCurrencyResponse, error) {
	fmt.Println("AttachCurrency")

	// Add product to the databased return the generated UUID
	product, err := s.productRepository.AttachCurrency(
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

func (s *economyServiceServer) DetachCurrency(ctx context.Context, req *v1.DetachCurrencyRequest) (*v1.DetachCurrencyResponse, error) {
	fmt.Println("DetachCurrency")

	// Add product to the databased return the generated UUID
	product, err := s.productRepository.DetachCurrency(
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

func (s *economyServiceServer) ListProductPrice(ctx context.Context, req *v1.ListProductPriceRequest) (*v1.ListProductPriceResponse, error) {
	fmt.Println("ListProductPrice")

	if req.GetProductId() == "" {
		return nil, status.Error(codes.InvalidArgument, "no product_id given")
	}

	prices, err := s.productRepository.ListPrice(ctx, req.GetProductId())
	if err != nil {
		return nil, err
	}

	return &v1.ListProductPriceResponse{
		Prices: prices,
	}, nil
}

func (s *economyServiceServer) BuyProduct(ctx context.Context, req *v1.BuyProductRequest) (*v1.BuyProductResponse, error) {
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
	product, err := s.productRepository.Get(ctx, req.GetProductId())
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
	payingStorage, err := s.storageRepository.Get(ctx, req.GetPayingStorageId())
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
		receivingStorage, err = s.storageRepository.Get(ctx, req.GetReceivingStorageId())
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

	_, err = s.productRepository.BuyProduct(
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
