package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	// "github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
)

func main() {
	// get configuration
	address := flag.String("server", "127.0.0.1:3000", "gRPC server in format host:port")
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
		ItemId: res1.Item.GetId(),
		Name:   "bierplezier",
	}

	res2, err := client.UpdateItem(ctx, &req2)
	if err != nil {
		log.Fatalf("UpdateItem failed: %v", err)
	}

	fmt.Println("UpdateItem result:")
	fmt.Println(proto.MarshalTextString(res2))

	// Call UpdateItem
	req6 := v1.CreateStorageRequest{
		PlayerId: "1337",
		Name:     "bank",
	}

	res6, err := client.CreateStorage(ctx, &req6)
	if err != nil {
		log.Fatalf("CreateStorage failed: %v", err)
	}

	fmt.Println("CreateStorage result:")
	fmt.Println(proto.MarshalTextString(res6))

	// Call GiveItem
	req7 := v1.GiveItemRequest{
		StorageId: res6.GetStorage().GetId(),
		ItemId:    res1.GetItem().GetId(),
	}

	res7, err := client.GiveItem(ctx, &req7)
	if err != nil {
		log.Fatalf("GiveItem failed: %v", err)
	}

	fmt.Println("GiveItem result:")
	fmt.Println(proto.MarshalTextString(res7))

	// Call GetStorage
	req8 := v1.GetStorageRequest{
		StorageId: res6.GetStorage().GetId(),
	}

	res8, err := client.GetStorage(ctx, &req8)
	if err != nil {
		log.Fatalf("GetStorage failed: %v", err)
	}

	fmt.Println("GetStorage result:")
	fmt.Println(proto.MarshalTextString(res8))
}
