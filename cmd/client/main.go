package main

import (
  "context"
  "flag"
  "log"
  "fmt"
  "time"

  // "github.com/golang/protobuf/ptypes"
  "google.golang.org/grpc"
  "github.com/golang/protobuf/proto"

  v1 "github.com/GameComponent/economy-service/pkg/api/v1"
)

const (
  // apiVersion is version of API is provided by server
  apiVersion = "v1"
)

func main() {
  // get configuration
  address := flag.String("server", "", "gRPC server in format host:port")
  flag.Parse()

  // Set up a connection to the server.
  conn, err := grpc.Dial(*address, grpc.WithInsecure())
  if err != nil {
    log.Fatalf("did not connect: %v", err)
  }
  defer conn.Close()

  client := v1.NewEconomyServiceClient(conn)

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  // Call CreateItem
  req1 := v1.CreateItemRequest{
    Api: apiVersion,
    Name: "kaasbaas",
  }

  res1, err := client.CreateItem(ctx, &req1)
  if err != nil {
    log.Fatalf("CreateItem failed: %v", err)
  }

  fmt.Println("CreateItem result:")
  fmt.Println(proto.MarshalTextString(res1))

  // Call UpdateItem
  req2 := v1.UpdateItemRequest{
    Api: apiVersion,
    ItemId: res1.Item.GetId(),
    Name: "bierplezier",
  }

  res2, err := client.UpdateItem(ctx, &req2)
  if err != nil {
    log.Fatalf("UpdateItem failed: %v", err)
  }

  fmt.Println("UpdateItem result:")
  fmt.Println(proto.MarshalTextString(res2))

  // Call UpdateItem
  req6 := v1.CreateStorageRequest{
    Api: apiVersion,
    PlayerId: "1337",
    StorageName: "bank",
  }

  res6, err := client.CreateStorage(ctx, &req6)
  if err != nil {
    log.Fatalf("CreateStorage failed: %v", err)
  }

  fmt.Println("CreateStorage result:")
  fmt.Println(proto.MarshalTextString(res6))

  // Call GiveItem
  req7 := v1.GiveItemRequest{
    Api: apiVersion,
    StorageId: res6.GetStorage().GetId(),
    ItemId: res1.GetItem().GetId(),
  }

  res7, err := client.GiveItem(ctx, &req7)
  if err != nil {
    log.Fatalf("GiveItem failed: %v", err)
  }

  fmt.Println("GiveItem result:")
  fmt.Println(proto.MarshalTextString(res7))

  // Call GetStorage
  req8 := v1.GetStorageRequest{
    Api: apiVersion,
    StorageId: res6.GetStorage().GetId(),
  }

  res8, err := client.GetStorage(ctx, &req8)
  if err != nil {
    log.Fatalf("GetStorage failed: %v", err)
  }

  fmt.Println("GetStorage result:")
  fmt.Println(proto.MarshalTextString(res8))

  // // Call UpdateItem3
  // req3 := v1.UpdateItemRequest{
  //   Api: apiVersion,
  //   ItemId: res1.Item.GetId(),
  //   Name: "bierplezier222",
  // }

  // res3, err := client.UpdateItem(ctx, &req3)
  // if err != nil {
  //   log.Fatalf("UpdateItem2 failed: %v", err)
  // }

  // fmt.Println("UpdateItem2 result:")
  // fmt.Println(proto.MarshalTextString(res3))

  // // Call CreateItem4
  // req4 := v1.CreateItemRequest{
  //   Api: apiVersion,
  //   Name: "dfgdfgfd777",
  // }

  // res4, err := client.CreateItem(ctx, &req4)
  // if err != nil {
  //   log.Fatalf("CreateItem failed: %v", err)
  // }

  // fmt.Println("CreateItem4 result:")
  // fmt.Println(proto.MarshalTextString(res4))

  // Call ListItems5
  // req5 := v1.ListItemsRequest{
  //   Api: apiVersion,
  // }

  // res5, err := client.ListItems(ctx, &req5)
  // if err != nil {
  //   log.Fatalf("ListItems failed: %v", err)
  // }

  // fmt.Println("ListItems result:")
  // fmt.Println(proto.MarshalTextString(res5))

  // Call Give
  // req1 := v1.GiveItemRequest{
  //   Api: apiVersion,
  //   StorageId: "storageid",
  //   ItemId: "testid",
  // }

  // res1, err := client.GiveItem(ctx, &req1)
  // if err != nil {
  //   log.Fatalf("Give failed: %v", err)
  // }

  // fmt.Println("Give result:")
  // fmt.Println(proto.MarshalTextString(res1))

  // Call GetInventory
  // req2 := v1.GetStorageRequest{
  //   Api: apiVersion,
  //   StorageId: "storageid",
  // }

  // res2, err := client.GetStorage(ctx, &req2)
  // if err != nil {
  //   log.Fatalf("Give failed: %v", err)
  // }

  // fmt.Println("GetInventory result:")
  // fmt.Println(proto.MarshalTextString(res2))

  // Call GetPlayer
  // req3 := v1.GetPlayerRequest{
  //   Api: apiVersion,
  //   PlayerId: "playerid",
  // }

  // res3, err := client.GetPlayer(ctx, &req3)
  // if err != nil {
  //   log.Fatalf("Give failed: %v", err)
  // }

  // fmt.Println("GetPlayer result:")
  // fmt.Println(proto.MarshalTextString(res3))
}
