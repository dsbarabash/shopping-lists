package main

import (
	"github.com/dsbarabash/shopping-lists/internal/handler"
	"github.com/dsbarabash/shopping-lists/internal/proto_api/pkg/grpc/v1/shopping_list_api"
	"github.com/dsbarabash/shopping-lists/internal/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	MongoDb, err := repository.ConnectMongoDb()
	if err != nil {
		log.Fatal(err)
	}
	RedisDB, err := repository.ConnectRedisDb()
	if err != nil {
		log.Fatal(err)
	}
	logWriter := repository.NewLogWriter(RedisDB)
	log.SetOutput(logWriter)

	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			handler.LoggingInterceptor,
		),
	)
	shopping_list_api.RegisterShoppingListServiceServer(s, &handler.GrpcServer{MongoDb: MongoDb})

	reflection.Register(s)

	log.Println("Server is running at :5001")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
