package entity

import (
	coreEnum "github.com/Rasikrr/learning_platform_core/enum"
	"github.com/google/uuid"
	"time"
)

//go:generate easyjson -all user.go

type User struct {
	ID          string               `json:"id"`
	Name        *string              `json:"name"`
	LastName    *string              `json:"last_name"`
	Email       string               `json:"email"`
	Password    string               `json:"password"`
	AccountRole coreEnum.AccountRole `json:"account_role"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	DeletedAt   *time.Time           `json:"deleted_at"`
}

func NewUser(email, password string) *User {
	return &User{
		ID:          uuid.New().String(),
		Email:       email,
		Password:    password,
		AccountRole: coreEnum.AccountRoleUser,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func NewAdminUser(email, password string) *User {
	return &User{
		ID:          uuid.New().String(),
		Email:       email,
		Password:    password,
		AccountRole: coreEnum.AccountRoleAdmin,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

type UpdateUserParams struct {
	ID       string  `json:"id"`
	Name     *string `json:"name"`
	LastName *string `json:"last_name"`
}
