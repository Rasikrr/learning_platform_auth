package grpc

import (
	"context"
	pb "github.com/Rasikrr/learning_platform_auth/pkg/api/proto/auth"
	"log"
)

func (s *Server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.EmptySuccessResponse, error) {
	log.Println(in)
	if err := s.authService.Register(ctx, in.GetEmail(), in.GetPassword(), in.GetPasswordConfirmation()); err != nil {
		return nil, err
	}
	return &pb.EmptySuccessResponse{}, nil
}

func (s *Server) ConfirmRegister(ctx context.Context, in *pb.ConfirmRegisterRequest) (*pb.AuthResponse, error) {
	tokens, err := s.authService.ConfirmRegister(ctx, in.GetEmail(), in.GetCode())
	if err != nil {
		return nil, err
	}
	return convertAuth(tokens), nil

}

func (s *Server) ConfirmAdminRegister(ctx context.Context, in *pb.ConfirmRegisterRequest) (*pb.AuthResponse, error) {
	tokens, err := s.authService.ConfirmAdminRegister(ctx, in.GetEmail(), in.GetCode())
	if err != nil {
		return nil, err
	}
	return convertAuth(tokens), nil
}

func (s *Server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.AuthResponse, error) {
	tokens, err := s.authService.Login(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {
		return nil, err
	}
	return convertAuth(tokens), nil
}

func (s *Server) ResetPassword(ctx context.Context, in *pb.ResetPasswordRequest) (*pb.EmptySuccessResponse, error) {
	if err := s.authService.ResetPassword(ctx, in.GetEmail(), in.GetPassword(), in.GetPasswordConfirmation()); err != nil {
		return nil, err
	}
	return &pb.EmptySuccessResponse{}, nil
}

func (s *Server) ConfirmResetPassword(ctx context.Context, in *pb.ConfirmResetPasswordRequest) (*pb.EmptySuccessResponse, error) {
	if err := s.authService.ConfirmResetPassword(ctx, in.GetEmail(), in.GetCode()); err != nil {
		return nil, err
	}
	return &pb.EmptySuccessResponse{}, nil
}

func (s *Server) CheckToken(ctx context.Context, in *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	ses, err := s.authService.CheckToken(ctx, in.GetToken())
	if err != nil {
		return nil, err
	}
	return &pb.CheckTokenResponse{
		Session: convertSession(ses),
	}, nil
}

func (s *Server) RefreshToken(ctx context.Context, in *pb.RefreshTokenRequest) (*pb.AuthResponse, error) {
	tokens, err := s.authService.RefreshToken(ctx, in.GetRefreshToken())
	if err != nil {
		return nil, err
	}
	return convertAuth(tokens), nil
}
