package main

import (
	grpc2 "github.com/dsbarabash/shopping-lists/internal/frontend/grpc"
	"github.com/dsbarabash/shopping-lists/internal/proto_api/pkg/grpc/v1/shopping_list_api"
	"github.com/dsbarabash/shopping-lists/internal/repository"
	"github.com/dsbarabash/shopping-lists/internal/repository/mongo"
	"github.com/dsbarabash/shopping-lists/internal/repository/redis"
	"github.com/dsbarabash/shopping-lists/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	MongoDb, err := mongo.ConnectMongoDb()
	if err != nil {
		log.Fatal(err)
	}
	RedisDB, err := redis.ConnectRedisDb()
	if err != nil {
		log.Fatal(err)
	}
	Service, err := service.NewListService(MongoDb)
	logWriter := repository.NewLogWriter(RedisDB)
	log.SetOutput(logWriter)

	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc2.LoggingInterceptor,
		),
	)
	shopping_list_api.RegisterShoppingListServiceServer(s, &grpc2.GrpcServer{Service: Service})

	reflection.Register(s)

	log.Println("Server is running at :5001")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
