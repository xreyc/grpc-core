package grpc

import (
	"log"

	authv1 "github.com/xreyc/grpc-core/internal/gen/go/auth/v1"
	"google.golang.org/grpc"
)

var AuthClient authv1.UserServiceClient

func InitGRPCClients() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to auth service: %v", err)
	}

	AuthClient = authv1.NewUserServiceClient(conn)
}
