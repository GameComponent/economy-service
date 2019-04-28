package shoprepository_test

import (
	"context"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	shoprepository "github.com/GameComponent/economy-service/pkg/repository/shop"
)

func TestGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := sqlmock.NewRows([]string{
		"shop.id",
		"shop.name",
		"shop.created_at",
		"shop.updated_at",
		"product.id",
		"product.name",
		"product.created_at",
		"product.updated_at",
		"item.id",
		"item.name",
		"item.metadata",
		"product_item.id",
		"product_item.amount",
	}).AddRow(
		"shop.id",
		"shop.name",
		time.Now(),
		time.Now(),
		"product.id",
		"product.name",
		time.Now(),
		time.Now(),
		"item.id",
		"item.name",
		"{}",
		"product_item.id",
		2,
	)
	mock.ExpectQuery("SELECT (.+)").WillReturnRows(columns)

	shopRepository := shoprepository.NewShopRepository(db)
	result, err := shopRepository.Get(context.Background(), "shop.id")
	if err != nil {
		t.Error(err)
	}

	if result == nil {
		t.Errorf("Result is nil")
	}

	if result.GetId() != "shop.id" {
		t.Errorf("result.GetId() does not match")
	}

	if result.GetName() != "shop.name" {
		t.Errorf("result.GetName() does not match")
	}

	products := result.GetProducts()
	if len(products) != 1 {
		t.Errorf("result.GetProducts() should return 1 product")
	}

	product := products[0]
	if product == nil {
		t.Errorf("product.GetProducts() first item should not be nil")
	}

	if product.GetId() != "product.id" {
		t.Errorf("product.GetId() does not match")
	}

	if product.GetName() != "product.name" {
		t.Errorf("product.GetName() does not match")
	}

	productItems := product.Items
	if len(productItems) != 1 {
		t.Errorf("product.GetProducts() should return 1 item")
	}

	productItem := productItems[0]
	if productItem == nil {
		t.Errorf("product.Items first item should not be nil")
	}

	if productItem.GetId() != "product_item.id" {
		t.Errorf("productItem.GetId() does not match")
	}

	if productItem.GetAmount() != 2 {
		t.Errorf("productItem.GetAmount() does not match")
	}

	item := productItem.Item
	if item == nil {
		t.Errorf("productItem.Item should not be nil")
	}

	if item.GetId() != "item.id" {
		t.Errorf("item.GetId() does not match")
	}

	if item.GetName() != "item.name" {
		t.Errorf("item.GetName() does not match")
	}
}

func TestGet2Products(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := sqlmock.NewRows([]string{
		"shop.id",
		"shop.name",
		"shop.created_at",
		"shop.updated_at",
		"product.id",
		"product.name",
		"product.created_at",
		"product.updated_at",
		"item.id",
		"item.name",
		"item.metadata",
		"product_item.id",
		"product_item.amount",
	}).AddRow(
		"shop.id",
		"shop.name",
		time.Now(),
		time.Now(),
		"product.id",
		"product.name",
		time.Now(),
		time.Now(),
		"item.id",
		"item.name",
		"{}",
		"product_item.id",
		2,
	).AddRow(
		"shop.id",
		"shop.name",
		time.Now(),
		time.Now(),
		"product.idb",
		"product.nameb",
		time.Now(),
		time.Now(),
		"item.idb",
		"item.nameb",
		"{}",
		"product_item.idb",
		2,
	)
	mock.ExpectQuery("SELECT (.+)").WillReturnRows(columns)

	shopRepository := shoprepository.NewShopRepository(db)
	result, err := shopRepository.Get(context.Background(), "shop.id")
	if err != nil {
		t.Error(err)
	}

	if result == nil {
		t.Errorf("Result is nil")
	}

	if result.GetId() != "shop.id" {
		t.Errorf("result.GetId() does not match")
	}

	if result.GetName() != "shop.name" {
		t.Errorf("result.GetName() does not match")
	}

	products := result.GetProducts()
	if len(products) != 2 {
		t.Errorf("result.GetProducts() should return 2 products")
	}

	product := products[0]
	productB := products[1]

	if product == nil {
		t.Errorf("product.GetProducts() first item should not be nil")
	}

	if product.GetId() != "product.id" {
		t.Errorf("product.GetId() does not match")
	}

	if productB.GetId() != "product.idb" {
		t.Errorf("productB.GetId() does not match")
	}

	if product.GetName() != "product.name" {
		t.Errorf("product.GetName() does not match")
	}

	if productB.GetName() != "product.nameb" {
		t.Errorf("productB.GetName() does not match")
	}
}

func TestGet2Items(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	columns := sqlmock.NewRows([]string{
		"shop.id",
		"shop.name",
		"shop.created_at",
		"shop.updated_at",
		"product.id",
		"product.name",
		"product.created_at",
		"product.updated_at",
		"item.id",
		"item.name",
		"item.metadata",
		"product_item.id",
		"product_item.amount",
	}).AddRow(
		"shop.id",
		"shop.name",
		time.Now(),
		time.Now(),
		"product.id",
		"product.name",
		time.Now(),
		time.Now(),
		"item.id",
		"item.name",
		"{}",
		"product_item.id",
		2,
	).AddRow(
		"shop.id",
		"shop.name",
		time.Now(),
		time.Now(),
		"product.id",
		"product.name",
		time.Now(),
		time.Now(),
		"item.idb",
		"item.nameb",
		"{}",
		"product_item.idb",
		4,
	)
	mock.ExpectQuery("SELECT (.+)").WillReturnRows(columns)

	shopRepository := shoprepository.NewShopRepository(db)
	result, err := shopRepository.Get(context.Background(), "shop.id")
	if err != nil {
		t.Error(err)
	}

	if result == nil {
		t.Errorf("Result is nil")
	}

	if result.GetId() != "shop.id" {
		t.Errorf("result.GetId() does not match")
	}

	if result.GetName() != "shop.name" {
		t.Errorf("result.GetName() does not match")
	}

	products := result.GetProducts()
	if len(products) != 1 {
		t.Errorf("result.GetProducts() should return 1 product")
	}

	product := products[0]
	if product == nil {
		t.Errorf("product.GetProducts() first item should not be nil")
	}

	if product.GetId() != "product.id" {
		t.Errorf("product.GetId() does not match")
	}

	if product.GetName() != "product.name" {
		t.Errorf("product.GetName() does not match")
	}

	productItems := product.Items
	if len(productItems) != 2 {
		t.Errorf("product.GetProducts() should return 2 items")
	}

	productItem := productItems[0]
	if productItem == nil {
		t.Errorf("product.Items first item should not be nil")
	}

	if productItem.GetId() != "product_item.id" {
		t.Errorf("productItem.GetId() does not match")
	}

	if productItem.GetAmount() != 2 {
		t.Errorf("productItem.GetAmount() does not match")
	}

	item := productItem.Item
	if item == nil {
		t.Errorf("productItem.Item should not be nil")
	}

	if item.GetId() != "item.id" {
		t.Errorf("item.GetId() does not match")
	}

	if item.GetName() != "item.name" {
		t.Errorf("item.GetName() does not match")
	}

	productItemb := productItems[1]
	if productItemb == nil {
		t.Errorf("product.Items second item should not be nil")
	}

	if productItemb.GetId() != "product_item.idb" {
		t.Errorf("productItemb.GetId() does not match")
	}

	if productItemb.GetAmount() != 4 {
		t.Errorf("productItemb.GetAmount() does not match")
	}

	itemb := productItemb.Item
	if itemb == nil {
		t.Errorf("productItemb.Item should not be nil")
	}

	if itemb.GetId() != "item.idb" {
		t.Errorf("itemb.GetId() does not match")
	}

	if itemb.GetName() != "item.nameb" {
		t.Errorf("itemb.GetName() does not match")
	}
}
