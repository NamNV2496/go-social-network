package grpc

import (
	"context"
	"fmt"

	"github.com/namnv2496/newsfeed-service/internal/configs"
	userv1 "github.com/namnv2496/newsfeed-service/internal/handler/generated/user_core/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductGRPCClient interface {
	GetFollowing(context.Context, string) ([]string, error)
}

type productGRPCClient struct {
	conn *grpc.ClientConn
}

func NewGRPCProductClient(
	config configs.GRPC,
) (ProductGRPCClient, error) {
	conn, err := grpc.NewClient(config.UserServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &productGRPCClient{
		conn: conn,
	}, nil
}

func (c productGRPCClient) GetFollowing(ctx context.Context, userId string) ([]string, error) {

	client := userv1.NewAccountServiceClient(c.conn)

	result, err := client.GetFollowing(context.Background(), &userv1.GetFollowingRequest{
		UserId: userId,
	})

	if err != nil {
		fmt.Println("Error: ", err)
	}
	// fmt.Println("result: ", result)
	return result.UserId, nil
}
