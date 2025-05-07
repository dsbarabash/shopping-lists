package main

import (
	"context"
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/proto_api/pkg/grpc/v1/shopping_list_api"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"math/rand"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.NewClient("localhost:5001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()
	unicIdx := rand.Intn(100)
	client := shopping_list_api.NewShoppingListServiceClient(conn)

	res3, err := client.CreateItem(context.Background(), &shopping_list_api.CreateItemRequest{Title: fmt.Sprintf("Item_title_%s", unicIdx), UserId: "343434", ShoppingListId: "567567567"})
	if err != nil {
		log.Fatalf("error create request: %s", err.Error())
	}
	fmt.Println("Create item ", res3.String())

	res2, err := client.GetItems(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("error get request: %s", err.Error())
	}
	fmt.Println("GetItems ", res2.String())

	res, err := client.GetItem(context.Background(), &shopping_list_api.GetItemRequest{Id: res2.Items[0].Id})
	if err != nil {
		log.Fatalf("error get request: %s", err.Error())
	}
	fmt.Println("GetItem ", res.String())

	res4, err := client.UpdateItem(context.Background(), &shopping_list_api.UpdateItemRequest{Id: res2.Items[0].Id, Title: fmt.Sprintf("Updated_Item_title_%s", unicIdx), UserId: "65656", ShoppingListId: "8989989"})
	if err != nil {
		log.Fatalf("error get request: %s", err.Error())
	}
	fmt.Println("UpdateItem ", res4.String())

	res5, err := client.DeleteItem(context.Background(), &shopping_list_api.DeleteItemRequest{Id: res2.Items[0].Id})
	if err != nil {
		log.Fatalf("error get request: %s", err.Error())
	}
	fmt.Println("DeleteItem ", res5.String())

	res8, err := client.CreateShoppingList(context.Background(), &shopping_list_api.CreateShoppingListRequest{Title: fmt.Sprintf("Sl_title_%s", unicIdx), UserId: "007", Items: []string{"234", "567"}})
	if err != nil {
		log.Fatalf("error get request: %s", err.Error())
	}
	fmt.Println("CreateShoppingList ", res8.String())

	res7, err := client.GetShoppingLists(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("error get request: %s", err.Error())
	}
	fmt.Println("GetShoppingLists ", res7.String())

	res6, err := client.GetShoppingList(context.Background(), &shopping_list_api.GetShoppingListRequest{Id: res7.ShoppingList[0].Id})
	if err != nil {
		log.Fatalf("error get request: %s", err.Error())
	}

	fmt.Println("GetShoppingList ", res6.String())

	res9, err := client.UpdateShoppingList(context.Background(), &shopping_list_api.UpdateShoppingListRequest{Id: res7.ShoppingList[0].Id, Title: fmt.Sprintf("Updated_Sl_title_%s", unicIdx), UserId: "123123", Items: []string{"234", "567", "875"}})
	if err != nil {
		log.Fatalf("error get request: %s", err.Error())
	}
	fmt.Println(res9.String())

	res10, err := client.DeleteShoppingList(context.Background(), &shopping_list_api.DeleteShoppingListRequest{Id: res7.ShoppingList[0].Id})
	if err != nil {
		log.Fatalf("error get request: %s", err.Error())
	}
	fmt.Println(res10.String())
}
