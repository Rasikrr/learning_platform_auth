package users

import (
	"github.com/Rasikrr/learning_platform_auth/internal/domain/entity"
	pb "github.com/Rasikrr/learning_platform_auth/pkg/deps/api/proto/users"
	coreEnum "github.com/Rasikrr/learning_platform_core/enum"
	"github.com/Rasikrr/learning_platform_core/grpc/converters"
)

func convert(res *pb.GetByEmailResponse) (*entity.User, error) {
	user := res.GetUser()
	if user == nil {
		return nil, nil
	}
	role, err := coreEnum.AccountRoleString(user.GetAccountRole())
	if err != nil {
		return nil, err
	}
	return &entity.User{
		ID:          user.Id,
		Name:        user.Name,
		LastName:    user.LastName,
		Email:       user.Email,
		Password:    user.Password,
		AccountRole: role,
		CreatedAt:   converters.ConvertToTime(user.CreatedAt),
		UpdatedAt:   converters.ConvertToTime(user.UpdatedAt),
		DeletedAt:   converters.ConvertToTimePtr(user.DeletedAt),
	}, nil
}

func convertUserToPb(user *entity.User) *pb.User {
	return &pb.User{
		Id:          user.ID,
		Name:        user.Name,
		LastName:    user.LastName,
		Email:       user.Email,
		Password:    user.Password,
		AccountRole: user.AccountRole.String(),
		CreatedAt:   converters.ConvertToTimestampPb(&user.CreatedAt),
		UpdatedAt:   converters.ConvertToTimestampPb(&user.UpdatedAt),
		DeletedAt:   converters.ConvertToTimestampPb(user.DeletedAt),
	}
}
