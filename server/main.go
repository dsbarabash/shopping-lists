package main

import (
	"github.com/dsbarabash/shopping-lists/internal/handler"
	"github.com/dsbarabash/shopping-lists/internal/proto_api/pkg/grpc/v1/shopping_list_api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	handler.FillSlices()
	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			handler.LoggingInterceptor,
		),
	)
	shopping_list_api.RegisterShoppingListServiceServer(s, &handler.Server{})

	reflection.Register(s)

	log.Println("Server is running at :5001")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
