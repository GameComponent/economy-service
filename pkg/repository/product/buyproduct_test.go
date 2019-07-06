package productrepository_test

import (
	"context"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	productrepository "github.com/GameComponent/economy-service/pkg/repository/product"
)

func TestBuyProductShouldStartTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT *")
	mock.ExpectCommit()

	productRepository := productrepository.NewProductRepository(db)
	product := v1.Product{}
	price := v1.Price{}
	receivingStorage := v1.Storage{}
	payingStorage := v1.Storage{}

	_, err = productRepository.BuyProduct(
		context.Background(),
		&product,
		&price,
		&receivingStorage,
		&payingStorage,
	)
	if err != nil {
		t.Error(err)
	}
}
