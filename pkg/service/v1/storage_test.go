package v1_test

import (
	"testing"
)

func TestGetExistingStorageItems(t *testing.T) {
	// db, _, err := sqlmock.New()
	// if err != nil {
	// 	t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	// }
	// defer db.Close()

	// storageRepository := storagerepository.NewStorageRepository(db)
	// itemRepository := itemrepository.NewItemRepository(db)

	// s := v1service.NewEconomyServiceServer(
	// 	db,
	// 	itemRepository,
	// 	nil,
	// 	nil,
	// 	storageRepository,
	// 	nil,
	// 	nil,
	// 	nil,
	// 	nil,
	// )

	// amount := v1.Amount{
	// 	MinAmount: 5,
	// 	MaxAmount: 5,
	// }

	// giveItemRequest := v1.GiveItemRequest{
	// 	ItemId: "item_id",
	// 	Amount: &amount,
	// }

	// result, err := s.GiveItem(
	// 	context.Background(),
	// 	&giveItemRequest,
	// )

	// if result.Amount != 5 {
	// 	t.Errorf("result.Amount should be 5")
	// }
}
