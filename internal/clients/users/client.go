package users

import (
	"context"
	"github.com/Rasikrr/learning_platform_auth/internal/domain/entity"
	pb "github.com/Rasikrr/learning_platform_auth/pkg/deps/api/proto/users"
	"github.com/Rasikrr/learning_platform_core/grpc"
	grpc2 "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type Client interface {
	Create(ctx context.Context, user *entity.User) error
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	ResetPassword(ctx context.Context, email, password string) error
}

type client struct {
	client pb.UsersClient
}

func NewClient(ctx context.Context, addr string) (Client, error) {
	conn, err := grpc.NewClient(
		ctx,
		addr,
		grpc2.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, err
	}
	return &client{
		client: pb.NewUsersClient(conn),
	}, nil
}

func (c *client) Create(ctx context.Context, user *entity.User) error {
	_, err := c.client.Create(ctx, &pb.CreateUserRequest{User: convertUserToPb(user)})
	return err
}

func (c *client) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	reply, err := c.client.GetByEmail(ctx, &pb.GetByEmailRequest{Email: email})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return convert(reply)
}

func (c *client) ResetPassword(ctx context.Context, email, password string) error {
	_, err := c.client.ResetPassword(ctx, &pb.ResetPasswordRequest{
		Email:        email,
		PasswordHash: password,
	})
	return err
}
