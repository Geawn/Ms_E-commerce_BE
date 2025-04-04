package utils

import (
	"context"
	"os"

	pb "github.com/Geawn/Ms_E-commerce_BE/account-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserServiceClient struct {
	client pb.UserServiceClient
	conn   *grpc.ClientConn
}

func NewUserServiceClient() (*UserServiceClient, error) {
	userServiceAddr := os.Getenv("USER_SERVICE_ADDR")
	if userServiceAddr == "" {
		userServiceAddr = "localhost:50051" // default address
	}

	conn, err := grpc.Dial(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewUserServiceClient(conn)
	return &UserServiceClient{
		client: client,
		conn:   conn,
	}, nil
}

func (c *UserServiceClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *UserServiceClient) CreateUser(ctx context.Context, email, firstName, lastName string) (uint64, error) {
	req := &pb.CreateUserRequest{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}

	resp, err := c.client.CreateUser(ctx, req)
	if err != nil {
		return 0, err
	}

	if !resp.Success {
		return 0, err
	}

	return resp.UserId, nil
}
