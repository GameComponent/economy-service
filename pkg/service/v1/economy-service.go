package v1

import (
  "context"
  "fmt"
  "log"
  // "errors"
  "database/sql"

  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"

  "github.com/google/uuid"
  // "github.com/golang/protobuf/ptypes/struct"
  v1 "github.com/GameComponent/economy-service/pkg/api/v1"
)

const (
  apiVersion = "v1"
)

// toDoServiceServer is implementation of v1.ToDoServiceServer proto interface
type economyServiceServer struct {
  db *sql.DB
}

// NewToDoServiceServer creates ToDo service
func NewEconomyServiceServer(db *sql.DB) v1.EconomyServiceServer {
  return &economyServiceServer{
    db,
  }
}

// checkAPI checks if the API version requested by client is supported by server
func (s *economyServiceServer) checkAPI(api string) error {
  // API version is "" means use current version of the service
  if len(api) > 0 {
    if apiVersion != api {
      return status.Errorf(codes.Unimplemented,
        "unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
    }
  }
  return nil
}

func (s *economyServiceServer) CreateItem(ctx context.Context, req *v1.CreateItemRequest) (*v1.CreateItemResponse, error) {
  fmt.Println("CreateItem");

  // check if the API version requested by client is supported by server
  if err := s.checkAPI(req.Api); err != nil {
    return nil, err
  }

  // Add item to the databased return the generated UUID
  lastInsertUuid := ""
  err := s.db.QueryRowContext(
    ctx,
    `INSERT INTO item(name) VALUES ($1) RETURNING id`,
    req.GetName(),
  ).Scan(&lastInsertUuid)

  if err != nil {
    return nil, err
  }

  // Generate the object based on the generated id and the requested name
  item := &v1.Item{
    Id: lastInsertUuid,
    Name: req.GetName(),
  }

  return &v1.CreateItemResponse{
    Api: apiVersion,
    Item: item,
  }, nil
}

func (s *economyServiceServer) UpdateItem(ctx context.Context, req *v1.UpdateItemRequest) (*v1.UpdateItemResponse, error) {
  fmt.Println("UpdateItem");

  // check if the API version requested by client is supported by server
  if err := s.checkAPI(req.Api); err != nil {
    return nil, err
  }

  result, err := s.db.ExecContext(
    ctx,
    `UPDATE item SET name = $1, data = $2 WHERE id = $3`,
    req.GetName(),
    `{"kaas":"baas"}`,
    req.GetItemId(),
  )

  if err != nil {
    return nil, err
  }

  fmt.Println("result")
  fmt.Println(result)

  item := &v1.Item{
    Id: req.GetItemId(),
    Name: req.GetName(),
  }

  return &v1.UpdateItemResponse{
    Api: apiVersion,
    Item: item,
  }, nil
}

func (s *economyServiceServer) ListItems(ctx context.Context, req *v1.ListItemsRequest) (*v1.ListItemsResponse, error) {
  fmt.Println("ListItems");

  // check if the API version requested by client is supported by server
  if err := s.checkAPI(req.Api); err != nil {
    return nil, err
  }

  // Query items from the database
  rows, err := s.db.QueryContext(ctx, "SELECT id, name FROM economy.item")
  if err != nil {
    log.Fatal(err)
  }
  defer rows.Close()

  // Unwrap rows into items
  items := []*v1.Item{}
  for rows.Next() {
    var item v1.Item
    err := rows.Scan(&item.Id, &item.Name)
    if err != nil {
        log.Fatalln(err)
    } 

    items = append(items, &item)
  }

  return &v1.ListItemsResponse{
    Api: apiVersion,
    Items: items,
  }, nil
}

func (s *economyServiceServer) GiveItem(ctx context.Context, req *v1.GiveItemRequest) (*v1.GiveItemResponse, error) {
  fmt.Println("GiveItem");

  // check if the API version requested by client is supported by server
  if err := s.checkAPI(req.Api); err != nil {
    return nil, err
  }

  // Add item to the databased return the generated UUID
  lastInsertUuid := ""
  err := s.db.QueryRowContext(
    ctx,
    `INSERT INTO storage_item(item_id, storage_id) VALUES ($1, $2) RETURNING id`,
    req.GetItemId(),
    req.GetStorageId(),
  ).Scan(&lastInsertUuid)

  if err != nil {
    return nil, err
  }

  item := &v1.Item{}

  err = s.db.QueryRowContext(
    ctx,
    `SELECT id, name FROM item WHERE id = $1`,
    req.GetItemId(),
  ).Scan(&item.Id, &item.Name)

  if err != nil {
    return nil, err
  }

  storageItem := &v1.StorageItem{
    Id: lastInsertUuid,
    Item: item,
    // Metadata: metadata,
  }

  return &v1.GiveItemResponse{
    Api: apiVersion,
    StorageId: req.GetStorageId(),
    Item: storageItem,
  }, nil
}

func (s *economyServiceServer) CreateStorage(ctx context.Context, req *v1.CreateStorageRequest) (*v1.CreateStorageResponse, error) {
  fmt.Println("CreateStorage");

  // check if the API version requested by client is supported by server
  if err := s.checkAPI(req.Api); err != nil {
    return nil, err
  }

  // Add item to the databased return the generated UUID
  lastInsertUuid := ""
  err := s.db.QueryRowContext(
    ctx,
    `INSERT INTO storage(player_id, name) VALUES ($1, $2) RETURNING id`,
    req.GetPlayerId(),
    req.GetName(),
  ).Scan(&lastInsertUuid)

  if err != nil {
    return nil, err
  }

  storage := &v1.Storage{
    Id: lastInsertUuid,
    PlayerId: req.GetPlayerId(),
    Name: req.GetName(),
  }

  return &v1.CreateStorageResponse{
    Api: apiVersion,
    Storage: storage,
  }, nil
}

func (s *economyServiceServer) GetStorage(ctx context.Context, req *v1.GetStorageRequest) (*v1.GetStorageResponse, error) {
  fmt.Println("GetStorage");

  // check if the API version requested by client is supported by server
  if err := s.checkAPI(req.Api); err != nil {
    return nil, err
  }

  // Check if the requested storage id is a valid UUID
  _, err := uuid.Parse(req.GetStorageId())
  if err != nil {
    return nil, err
  }

  rows, err := s.db.QueryContext(
    ctx,
    `SELECT
      storage.id as storageId,
      storage.name as storageName,
      storage.data as storageData,
      storage.player_id as playerId,
      storage_item.id as storageItemId,
      storage_item.data as storageItemData,
      item.id as itemId,
      item.name as itemName,
      item.data as itemData
    FROM storage 
    INNER JOIN storage_item on (storage.id = storage_item.storage_id)
    INNER JOIN item on (storage_item.item_id = item.id)
    WHERE storage.id = $1`,
    req.GetStorageId(),
  )

  if err != nil {
    return nil, err
  }

  storageItems := []*v1.StorageItem{}

  type row struct {
    StorageId string
    StorageName string
    StorageData string
    PlayerId string
    StorageItemId string
    StorageItemItemData string
    ItemId string
    ItemName string
    ItemData string
  }

  var r row
  for rows.Next() {
    err = rows.Scan(
      &r.StorageId,
      &r.StorageName,
      &r.StorageData,
      &r.PlayerId,
      &r.StorageItemId,
      &r.StorageItemItemData,
      &r.ItemId,
      &r.ItemName,
      &r.ItemData,
    )
    if err != nil {
      log.Fatal(err)
    }

    fmt.Printf("rowzz")
    fmt.Printf("%+v\n", r)

    item := &v1.Item {
      Id: r.ItemId,
      Name: r.ItemName,
    }

    storageItem := &v1.StorageItem {
      Id: r.StorageItemId,
      Item: item,
    }

    storageItems = append(storageItems, storageItem)
  }

  storage := &v1.Storage {
    Id: r.StorageId,
    PlayerId: r.PlayerId,
    Name: r.StorageName,
    Items: storageItems,
  }
  
  return &v1.GetStorageResponse{
    Api: apiVersion,
    Storage: storage,
  }, nil
}

func (s *economyServiceServer) GetPlayer(ctx context.Context, req *v1.GetPlayerRequest) (*v1.GetPlayerResponse, error) {
  fmt.Println("GetPlayer");

  // check if the API version requested by client is supported by server
  if err := s.checkAPI(req.Api); err != nil {
    return nil, err
  }

  rows, err := s.db.QueryContext(
    ctx,
    `SELECT id, name
    FROM storage 
    WHERE player_id = $1`,
    req.GetPlayerId(),
  )

  if err != nil {
    return nil, err
  }

  storageItems := []*v1.StorageBase{}

  for rows.Next() {
    storage := &v1.StorageBase{}
    
    err = rows.Scan(
      &storage.Id,
      &storage.Name,
    )
    if err != nil {
      log.Fatal(err)
    }

    storageItems = append(storageItems, storage)
  }

  return &v1.GetPlayerResponse{
    Api: apiVersion,
    PlayerId: req.GetPlayerId(),
    Storages: storageItems,
  }, nil
}