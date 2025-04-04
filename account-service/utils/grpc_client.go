package utils

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Geawn/Ms_E-commerce_BE/account-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserServiceClient struct {
	client proto.UserServiceClient
	conn   *grpc.ClientConn
}

func NewUserServiceClient() (*UserServiceClient, error) {
	userServiceAddr := os.Getenv("USER_SERVICE_ADDR")
	if userServiceAddr == "" {
		userServiceAddr = "localhost:50052" // default address
	}

	log.Printf("Connecting to user-service at %s", userServiceAddr)
	conn, err := grpc.Dial(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to user-service: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to user-service")
	client := proto.NewUserServiceClient(conn)
	return &UserServiceClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *UserServiceClient) Close() {
	if c.conn != nil {
		c.conn.Close()
		log.Println("Closed connection to user-service")
	}
}

func (c *UserServiceClient) CreateUser(ctx context.Context, email, firstName, lastName string) (uint64, error) {
	req := &proto.CreateUserRequest{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}

	log.Printf("Sending gRPC request to create user: %+v", req)
	resp, err := c.client.CreateUser(ctx, req)
	if err != nil {
		log.Printf("gRPC call failed: %v", err)
		return 0, err
	}

	log.Printf("Received gRPC response: %+v", resp)
	if !resp.Success {
		log.Printf("Failed to create user in user-service: %s", resp.Error)
		return 0, fmt.Errorf("failed to create user in user service: %s", resp.Error)
	}

	log.Printf("Successfully created user with ID: %d", resp.UserId)
	return resp.UserId, nil
}
