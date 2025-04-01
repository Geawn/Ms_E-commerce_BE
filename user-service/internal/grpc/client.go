package grpc

import (
	"context"
	"fmt"

	pb "github.com/Geawn/Ms_E-commerce_BE/user-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.UserServiceClient
}

func NewClient(serverAddr string) (*Client, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}

	return &Client{
		conn:   conn,
		client: pb.NewUserServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetCurrentUser(ctx context.Context, userID string) (*pb.User, error) {
	resp, err := c.client.GetCurrentUser(ctx, &pb.GetCurrentUserRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, err
	}
	return resp.User, nil
}
