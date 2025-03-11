package grpc

import (
	"github.com/Rasikrr/learning_platform_auth/internal/domain/entity"
	pb "github.com/Rasikrr/learning_platform_auth/pkg/api/proto/auth"
	"github.com/Rasikrr/learning_platform_core/http/session"
)

func convertAuth(auth *entity.Auth) *pb.AuthResponse {
	return &pb.AuthResponse{
		AccessToken:  auth.AccessToken,
		RefreshToken: auth.RefreshToken,
	}
}

func convertSession(ses *session.Session) *pb.Session {
	claims := make(map[string]string)
	if ses.Claims() != nil {
		for k, v := range ses.Claims() {
			val, ok := v.(string)
			if !ok {
				continue
			}
			claims[k] = val
		}
	}
	return &pb.Session{
		UserId: ses.UserID(),
		Email:  ses.Email(),
		Role:   ses.AccountRole().String(),
		Claims: claims,
	}
}
