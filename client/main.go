package main

import (
	"context"
	"fmt"
	"github.com/dsbarabash/shopping-lists/internal/proto_api/pkg/grpc/v1/shopping_list_api"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.NewClient("localhost:5001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	client := shopping_list_api.NewShoppingListServiceClient(conn)

	res2, err := client.GetItems(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("error get request: %s", err.Error())
	}
	fmt.Println(res2.String())

	res, err := client.GetItem(context.Background(), &shopping_list_api.GetItemRequest{Id: res2.Items[0].Id})
	if err != nil {
		log.Fatalf("error get request: %s", err.Error())
	}

	fmt.Println(res.String())

	res3, err := client.CreateItem(context.Background(), &shopping_list_api.CreateItemRequest{Title: "Test_title1", UserId: "123123", ShoppingListId: "567567567"})
	if err != nil {
		log.Fatalf("error get request: %s", err.Error())
	}
	fmt.Println(res3.String())

	res4, err := client.UpdateItem(context.Background(), &shopping_list_api.UpdateItemRequest{Id: res2.Items[0].Id, Title: "Test_title3", UserId: "65656", ShoppingListId: "8989989"})
	if err != nil {
		log.Fatalf("error get request: %s", err.Error())
	}
	fmt.Println(res4.String())

	res5, err := client.DeleteItem(context.Background(), &shopping_list_api.DeleteItemRequest{Id: res2.Items[0].Id})
	if err != nil {
		log.Fatalf("error get request: %s", err.Error())
	}
	fmt.Println(res5.String())

	res7, err := client.GetShoppingLists(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("error get request: %s", err.Error())
	}
	fmt.Println(res7.String())

	res6, err := client.GetShoppingList(context.Background(), &shopping_list_api.GetShoppingListRequest{Id: res7.ShoppingList[0].Id})
	if err != nil {
		log.Fatalf("error get request: %s", err.Error())
	}

	fmt.Println(res6.String())

	res8, err := client.CreateShoppingList(context.Background(), &shopping_list_api.CreateShoppingListRequest{Title: "Test_title1", UserId: "0000", Items: []string{"234", "567"}})
	if err != nil {
		log.Fatalf("error get request: %s", err.Error())
	}
	fmt.Println(res8.String())

	res9, err := client.UpdateShoppingList(context.Background(), &shopping_list_api.UpdateShoppingListRequest{Id: res7.ShoppingList[0].Id, Title: "Test_title3", UserId: "123123", Items: []string{"234", "567", "875"}})
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
