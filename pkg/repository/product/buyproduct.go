package productrepository

import (
	"context"
	"database/sql"
	"fmt"
	"math"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
)

// BuyProduct buys a product
func (r *ProductRepository) BuyProduct(ctx context.Context, product *v1.Product, price *v1.Price, receivingStorage *v1.Storage, payingStorage *v1.Storage) (*v1.Product, error) {
	options := sql.TxOptions{
		ReadOnly: false,
	}

	// Start a transaction
	tx, err := r.db.BeginTx(ctx, &options)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Take the Currencies from the Storage
	err = takeCurrenciesFromStorage(ctx, tx, price.Currencies, payingStorage)
	if err != nil {
		return nil, err
	}

	// Take the Items from the Storage
	err = takeItemsFromStorage(ctx, tx, price.Items, payingStorage)
	if err != nil {
		return nil, err
	}

	// Give the Currencies from the Storage
	err = giveCurrenciesToStorage(ctx, tx, product.Currencies, receivingStorage)
	if err != nil {
		return nil, err
	}

	// Give items to the storage
	err = giveItemsToStorage(ctx, tx, product.Items, receivingStorage)
	if err != nil {
		return nil, err
	}

	// Commit all changes to the database
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("wat")
}

func takeCurrenciesFromStorage(ctx context.Context, tx *sql.Tx, priceCurrencies []*v1.PriceCurrency, storage *v1.Storage) error {
	// Get a map of paying StorageCurrencies
	payingStorageCurrencies := storage.Currencies
	payingStorageCurrenciesMap := map[string]*v1.StorageCurrency{}
	for _, payingStorageCurrency := range payingStorageCurrencies {
		payingStorageCurrenciesMap[payingStorageCurrency.Currency.Id] = payingStorageCurrency
	}

	// Take the Currencies from the Storage
	for _, priceCurrency := range priceCurrencies {
		// Get the StorageCurrency Id
		storageCurrency := payingStorageCurrenciesMap[priceCurrency.Currency.Id]

		_, err := tx.ExecContext(
			ctx,
			`
				UPDATE storage_currency
				SET amount = storage_currency.amount - $1
				WHERE storage_currency.id = $2
				AND storage_currency.amount = $3
			`,
			priceCurrency.Amount,
			storageCurrency.Id,
			storageCurrency.Amount,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func takeItemsFromStorage(ctx context.Context, tx *sql.Tx, priceItems []*v1.PriceItem, storage *v1.Storage) error {
	for _, priceItem := range priceItems {
		// Take the stackable items
		err := takeStackableItemFromStorage(ctx, tx, priceItem, storage)
		if err != nil {
			return err
		}

		// Take the unstackable items
		err = takeUnstackableItemFromStorage(ctx, tx, priceItem, storage)
		if err != nil {
			return err
		}
	}

	return nil
}

func takeStackableItemFromStorage(ctx context.Context, tx *sql.Tx, priceItem *v1.PriceItem, storage *v1.Storage) error {
	if priceItem.Item.Stackable == false {
		return nil
	}

	rows, err := tx.QueryContext(
		ctx,
		`
			SELECT id, amount
			FROM storage_item
			WHERE storage_id = $1
			AND item_id = $2
			ORDER BY amount DESC
		`,
		storage.Id,
		priceItem.Item.Id,
	)
	if err != nil {
		return err
	}

	// Check the amount of StorageItems
	amounts := []*v1.StorageItem{}
	for rows.Next() {
		amount := v1.StorageItem{}

		err := rows.Scan(
			&amount.Id,
			&amount.Amount,
		)
		if err != nil {
			return err
		}

		amounts = append(amounts, &amount)
	}

	// Check to make sure there are enough StorageItems in the Storage
	total := int64(0)
	for _, amount := range amounts {
		total = total + amount.Amount
	}
	if total < priceItem.Amount {
		return fmt.Errorf("not enough items in storage")
	}

	remainder := priceItem.Amount
	for _, amount := range amounts {
		if remainder == 0 {
			continue
		}

		// Calculate the amount to remove
		amountToRemove := amount.Amount
		if amountToRemove > remainder {
			amountToRemove = remainder
		}

		// Calculate the new remainder
		remainder = remainder - amountToRemove

		// Remove the entire stack
		if amountToRemove == amount.Amount {
			_, err := tx.ExecContext(
				ctx,
				`
					DELETE storage_item
					WHERE id = $1
				`,
				amount.Id,
			)
			if err != nil {
				return err
			}

			continue
		}

		// Remove some amount of a stack
		_, err := tx.ExecContext(
			ctx,
			`
				UPDATE storage_item
				SET amount = amount - $1
				WHERE id = $2
			`,
			amountToRemove,
			amount.Id,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func takeUnstackableItemFromStorage(ctx context.Context, tx *sql.Tx, priceItem *v1.PriceItem, storage *v1.Storage) error {
	if priceItem.Item.Stackable == true {
		return nil
	}

	// Fetch the amount of items in the storage
	var itemAmount int64
	tx.QueryRowContext(
		ctx,
		`
			SELECT COUNT(id)
			FROM storage_item
			WHERE storage_id = $1
			AND item_id = $2
		`,
		storage.Id,
		priceItem.Item.Id,
	).Scan(&itemAmount)

	// Check if there enough items in the Storage
	if itemAmount < priceItem.Amount {
		return fmt.Errorf("not enough items in storage")
	}

	// Delete the items from the storage
	_, err := tx.ExecContext(
		ctx,
		`
			DELETE FROM storage_item
			WHERE storage_id = $1
			AND item_id = $2
			LIMIT $3
		`,
		storage.Id,
		priceItem.Item.Id,
		priceItem.Amount,
	)
	if err != nil {
		return err
	}

	return nil
}

func giveCurrenciesToStorage(ctx context.Context, tx *sql.Tx, productCurrencies []*v1.ProductCurrency, storage *v1.Storage) error {
	for _, productCurrency := range productCurrencies {
		_, err := tx.ExecContext(
			ctx,
			`
				INSERT INTO storage_currency(currency_id, storage_id, amount)
				VALUES($1, $2, $3)
				ON CONFLICT(currency_id,storage_id) DO UPDATE
				SET amount = storage_currency.amount + EXCLUDED.amount
			`,
			productCurrency.Currency.Id,
			storage.Id,
			productCurrency.Amount,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func giveItemsToStorage(ctx context.Context, tx *sql.Tx, productItems []*v1.ProductItem, storage *v1.Storage) error {
	for _, productItem := range productItems {
		// Check if we can just insert a new stack
		giveDefaultStackableItemToStorage(ctx, tx, productItem, storage)

		// Fill existing stacks
		giveFillableStackableItemToStorage(ctx, tx, productItem, storage)

		// Add non-stackable items
		giveUnstackableItemToStorage(ctx, tx, productItem, storage)
	}

	return nil
}

func giveUnstackableItemToStorage(ctx context.Context, tx *sql.Tx, productItem *v1.ProductItem, storage *v1.Storage) error {
	if productItem.Item.Stackable == true {
		return nil
	}

	loops := int(productItem.Amount)
	for i := 0; i < loops; i++ {
		if err := addItemToStorage(ctx, tx, productItem.Item.Id, storage.Id, 1); err != nil {
			return err
		}
	}

	return nil
}

func giveDefaultStackableItemToStorage(ctx context.Context, tx *sql.Tx, productItem *v1.ProductItem, storage *v1.Storage) error {
	if !productItem.Item.Stackable {
		return nil
	}

	if productItem.Item.StackBalancingMethod != v1.StackBalancingMethod_UNBALANCED_CREATE_NEW_STACKS &&
		productItem.Item.StackBalancingMethod != v1.StackBalancingMethod_DEFAULT {
		return nil
	}

	err := addItemToStorage(
		ctx,
		tx,
		productItem.Item.Id,
		storage.Id,
		productItem.Amount,
	)
	if err != nil {
		return err
	}

	return nil
}

func giveFillableStackableItemToStorage(ctx context.Context, tx *sql.Tx, productItem *v1.ProductItem, storage *v1.Storage) error {
	if productItem.Item.Stackable == false {
		return nil
	}

	if productItem.Item.StackBalancingMethod != v1.StackBalancingMethod_UNBALANCED_FILL_EXISTING_STACKS &&
		productItem.Item.StackBalancingMethod != v1.StackBalancingMethod_BALANCED_FILL_EXISTING_STACKS {
		return nil
	}

	remainder := productItem.Amount

	// Fill existing stacks that are not yet full
	remainder, err := fillExistingStacks(ctx, tx, productItem, storage, remainder)
	if err != nil {
		return err
	}

	// Create new stack(s) for the remainder
	err = addStackableItemToStorage(ctx, tx, productItem, storage, remainder)
	if err != nil {
		return err
	}

	return nil
}

func fillExistingStacks(ctx context.Context, tx *sql.Tx, productItem *v1.ProductItem, storage *v1.Storage, remainder int64) (int64, error) {
	for _, receivingStorageItem := range storage.Items {
		if productItem.Item.Id != receivingStorageItem.Item.Id {
			continue
		}

		if productItem.Item.StackMaxAmount > 0 && receivingStorageItem.Amount >= productItem.Item.StackMaxAmount {
			continue
		}

		// Calculate how many space is left in this StorageItem
		available := remainder
		if productItem.Item.StackMaxAmount > 0 {
			available = productItem.Item.StackMaxAmount - receivingStorageItem.Amount
		}

		remainder = remainder - available
		_, err := tx.ExecContext(
			ctx,
			`
				UPDATE storage_item
				SET amount = amount + $1
				WHERE id = $2
			`,
			available,
			receivingStorageItem.Id,
		)
		if err != nil {
			return 0, err
		}
	}

	return remainder, nil
}

func addItemToStorage(ctx context.Context, tx *sql.Tx, itemID string, storageID string, amount int64) error {
	_, err := tx.ExecContext(
		ctx,
		`INSERT INTO storage_item(item_id, storage_id, amount) VALUES ($1, $2, $3)`,
		itemID,
		storageID,
		amount,
	)

	return err
}

func addStackableItemToStorage(ctx context.Context, tx *sql.Tx, productItem *v1.ProductItem, storage *v1.Storage, amount int64) error {

	// Create fully filled stacks
	fullStacksToCreate := int(math.Floor(float64(amount) / float64(productItem.Item.StackMaxAmount)))
	for i := 0; i < fullStacksToCreate; i++ {
		err := addItemToStorage(
			ctx,
			tx,
			productItem.Item.Id,
			storage.Id,
			productItem.Item.StackMaxAmount,
		)
		if err != nil {
			return err
		}
	}

	// Create partial filled stacks
	partialStackToCreate := amount % productItem.Item.StackMaxAmount
	if partialStackToCreate > 0 {
		err := addItemToStorage(
			ctx,
			tx,
			productItem.Item.Id,
			storage.Id,
			partialStackToCreate,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
