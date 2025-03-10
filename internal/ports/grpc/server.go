package grpc

import (
	authS "github.com/Rasikrr/learning_platform_auth/internal/services/auth"
	pb "github.com/Rasikrr/learning_platform_auth/pkg/api/proto/auth"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedAuthServer
	authService authS.Service
}

func NewServer(
	server *grpc.Server,
	authService authS.Service) *Server {
	srv := &Server{
		authService: authService,
	}
	pb.RegisterAuthServer(server, srv)
	return srv
}
