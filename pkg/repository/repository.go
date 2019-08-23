package repository

import (
	"context"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
)

// AccountRepository interface
type AccountRepository interface {
	Create(ctx context.Context, email string, password string) (*v1.Account, error)
	Update(ctx context.Context, accountID string, password string) (*v1.Account, error)
	Get(ctx context.Context, accountID string) (*v1.Account, error)
	GetByEmail(ctx context.Context, email string) (*v1.Account, error)
	AssignPermission(ctx context.Context, accountID string, permission string) (*v1.Account, error)
	RevokePermission(ctx context.Context, accountID string, permission string) (*v1.Account, error)
}

// ConfigRepository interface
type ConfigRepository interface {
	Get(ctx context.Context, key string) (*v1.Config, error)
	Set(ctx context.Context, key string, value string) (*v1.Config, error)
	List(ctx context.Context, limit int32, offset int32) ([]*v1.Config, int32, error)
}

// CurrencyRepository interface
type CurrencyRepository interface {
	Create(ctx context.Context, name string, shortName string, symbol string) (*v1.Currency, error)
	Update(ctx context.Context, currencyID string, name string, shortName string, symbol string) (*v1.Currency, error)
	Get(ctx context.Context, currencyID string) (*v1.Currency, error)
	List(ctx context.Context, limit int32, offset int32) ([]*v1.Currency, int32, error)
}

// ItemRepository interface
type ItemRepository interface {
	Create(ctx context.Context, name string, stackable bool, stackMaxAmount int64, stackBalancingMethod int64, metadata string) (*v1.Item, error)
	Get(ctx context.Context, itemID string) (*v1.Item, error)
	Update(ctx context.Context, itemID string, name string, metadata string) (*v1.Item, error)
	List(ctx context.Context, limit int32, offset int32) ([]*v1.Item, int32, error)
	Search(ctx context.Context, query string, limit int32, offset int32) ([]*v1.Item, int32, error)
}

// PlayerRepository interface
type PlayerRepository interface {
	Create(ctx context.Context, playerID string, name string, metadata string) (*v1.Player, error)
	Update(ctx context.Context, playerID string, name string, metadata string) (*v1.Player, error)
	Get(ctx context.Context, playerID string) (*v1.Player, error)
	List(ctx context.Context, limit int32, offset int32) ([]*v1.Player, int32, error)
	Search(ctx context.Context, query string, limit int32, offset int32) ([]*v1.Player, int32, error)
}

// PriceRepository interface
type PriceRepository interface {
	Get(ctx context.Context, priceID string) (*v1.Price, error)
	Create(ctx context.Context, productID string) (*v1.Price, error)
	Delete(ctx context.Context, priceID string) (bool, error)
	AttachPriceCurrency(ctx context.Context, priceID string, currencyID string, amount int64) (*v1.Price, error)
	DetachPriceCurrency(ctx context.Context, priceCurrencyID string) (*v1.Price, error)
	AttachPriceItem(ctx context.Context, priceID string, itemID string, amount int64) (*v1.Price, error)
	DetachPriceItem(ctx context.Context, priceItemID string) (*v1.Price, error)
}

// ProductRepository interface
type ProductRepository interface {
	Create(ctx context.Context, name string) (*v1.Product, error)
	Get(ctx context.Context, productID string) (*v1.Product, error)
	Update(ctx context.Context, productID string, name string) (*v1.Product, error)
	List(ctx context.Context, limit int32, offset int32) ([]*v1.Product, int32, error)
	Search(ctx context.Context, query string, limit int32, offset int32) ([]*v1.Product, int32, error)
	AttachItem(ctx context.Context, productID string, itemID string, amount int64) (*v1.Product, error)
	DetachItem(ctx context.Context, productItemID string) (*v1.Product, error)
	AttachCurrency(ctx context.Context, productID string, currencyID string, amount int64) (*v1.Product, error)
	DetachCurrency(ctx context.Context, productCurrencyID string) (*v1.Product, error)
	BuyProduct(ctx context.Context, product *v1.Product, price *v1.Price, receivingStorage *v1.Storage, payingStorage *v1.Storage) (*v1.Product, error)
	ListPrice(ctx context.Context, productID string) ([]*v1.Price, error)
}

// ShopRepository interface
type ShopRepository interface {
	Get(ctx context.Context, shopID string) (*v1.Shop, error)
	Create(ctx context.Context, name string, metadata string) (*v1.Shop, error)
	Update(ctx context.Context, shopID string, name string, metadata string) (*v1.Shop, error)
	List(ctx context.Context, limit int32, offset int32) ([]*v1.Shop, int32, error)
	AttachProduct(ctx context.Context, shopID string, productID string) (*v1.Shop, error)
	DetachProduct(ctx context.Context, shopProductID string) (*v1.Shop, error)
}

// StorageRepository interface
type StorageRepository interface {
	Create(ctx context.Context, playerID string, name string, metadata string) (*v1.Storage, error)
	Update(ctx context.Context, storageID string, name string, metadata string) (*v1.Storage, error)
	Get(ctx context.Context, storageID string) (*v1.Storage, error)
	GiveItem(ctx context.Context, storageID string, itemID string, amount int64, metadata string) (*string, error)
	IncreaseItemAmount(ctx context.Context, storageItemID string, amount int64) error
	GiveCurrency(ctx context.Context, storageID string, currencyID string, amount int64) (*v1.StorageCurrency, error)
	List(ctx context.Context, limit int32, offset int32) ([]*v1.Storage, int32, error)
	SplitStack(ctx context.Context, storageItemID string, amounts []int64) (*v1.Storage, error)
	MergeStack(ctx context.Context, toStorageItemID string, fromStorageItemID string) (*v1.Storage, error)
}
