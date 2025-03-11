package auth

import (
	"context"
	"errors"
	authC "github.com/Rasikrr/learning_platform_auth/internal/cache/auth"
	usersC "github.com/Rasikrr/learning_platform_auth/internal/clients/users"
	"github.com/Rasikrr/learning_platform_auth/internal/domain/entity"
	"github.com/Rasikrr/learning_platform_auth/internal/util"
	coreJwt "github.com/Rasikrr/learning_platform_core/http/jwt"
	"github.com/Rasikrr/learning_platform_core/http/session"
	"log"
	"math/rand/v2"
	"strings"
	"time"
	"unicode"
)

const (
	codeLength        = 4
	passwordMinLength = 8
)

type Service interface {
	Register(ctx context.Context, email, password, passwordConfirmation string) error
	ConfirmRegister(ctx context.Context, email, code string) (*entity.Auth, error)
	ConfirmAdminRegister(ctx context.Context, email, code string) (*entity.Auth, error)
	Login(ctx context.Context, email, password string) (*entity.Auth, error)
	ResetPassword(ctx context.Context, email, password, passwordConfirmation string) error
	ConfirmResetPassword(ctx context.Context, email, code string) error
	CheckToken(ctx context.Context, token string) (*session.Session, error)
	RefreshToken(ctx context.Context, token string) (*entity.Auth, error)
}

type service struct {
	accessTTL   time.Duration
	refreshTTL  time.Duration
	cache       authC.Cache
	usersClient usersC.Client
}

func NewService(
	accessTTL, refreshTTL time.Duration,
	authCache authC.Cache,
	usersClient usersC.Client,
) Service {
	return &service{
		accessTTL:   accessTTL,
		refreshTTL:  refreshTTL,
		cache:       authCache,
		usersClient: usersClient,
	}
}

func (s *service) Login(ctx context.Context, email, password string) (*entity.Auth, error) {
	user, err := s.usersClient.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	if err := util.CheckPassword(user.Password, password); err != nil {
		return nil, err
	}
	return generateTokens(user, s.accessTTL, s.refreshTTL)
}

func (s *service) Register(ctx context.Context, email, password, passwordConfirm string) error {
	if password != passwordConfirm {
		return errors.New("passwords do not match")
	}
	user, err := s.usersClient.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user != nil {
		return errors.New("user already exists")
	}

	if err := validatePassword(password); err != nil {
		log.Println("password: ", password)
		return err
	}

	code := generateCode()
	//if err := s.mailClient.Send(ctx, []string{email}, "registration", code); err != nil {
	//	return err
	//}
	if err := s.cache.SetCode(ctx, email, code); err != nil {
		return err
	}
	passHash, err := util.Hash(password)
	if err != nil {
		return err
	}
	if err := s.cache.SetPasswordHash(ctx, email, passHash); err != nil {
		return err
	}
	return nil
}

func (s *service) ResetPassword(ctx context.Context, email, password, passwordConfirm string) error {
	if password != passwordConfirm {
		return errors.New("passwords do not match")
	}
	user, err := s.usersClient.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}
	err = validatePassword(password)
	if err != nil {
		return err
	}
	code := generateCode()
	//err = s.mailClient.Send(ctx, []string{email}, "reset_password", code)
	if err != nil {
		return err
	}
	if err := s.cache.SetCode(ctx, email, code); err != nil {
		return err
	}
	passHash, err := util.Hash(password)
	if err != nil {
		return err
	}
	return s.cache.SetPasswordHash(ctx, email, passHash)
}

func (s *service) ConfirmResetPassword(ctx context.Context, email, code string) error {
	codeFromCache, err := s.cache.GetCode(ctx, email)
	if err != nil {
		return err
	}
	if codeFromCache != code {
		return errors.New("code is invalid")
	}
	passHash, err := s.cache.GetPasswordHash(ctx, email)
	if err != nil {
		return err
	}
	if err = s.usersClient.ResetPassword(ctx, email, passHash); err != nil {
		return err
	}
	return nil
}

func (s *service) ConfirmRegister(ctx context.Context, email, code string) (*entity.Auth, error) {
	codeFromCache, err := s.cache.GetCode(ctx, email)
	if err != nil {
		return nil, err
	}
	if codeFromCache != code {
		return nil, errors.New("code is invalid")
	}
	passHash, err := s.cache.GetPasswordHash(ctx, email)
	if err != nil {
		return nil, err
	}
	user := entity.NewUser(email, passHash)
	if err := s.usersClient.Create(ctx, user); err != nil {
		return nil, err
	}

	return generateTokens(user, s.accessTTL, s.refreshTTL)
}

func (s *service) CheckToken(_ context.Context, token string) (*session.Session, error) {
	ses, isRefresh, err := coreJwt.ParseJwt(token)
	if err != nil {
		return nil, err
	}
	if isRefresh {
		return nil, errors.New("access token expected")
	}
	return ses, nil
}

func (s *service) RefreshToken(ctx context.Context, token string) (*entity.Auth, error) {
	ses, isRefresh, err := coreJwt.ParseJwt(token)
	if err != nil {
		return nil, err
	}
	if !isRefresh {
		return nil, errors.New("refresh token expected")
	}
	email := ses.Email()
	if email == "" {
		return nil, errors.New("email is empty")
	}
	user, err := s.usersClient.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return generateTokens(user, s.accessTTL, s.refreshTTL)
}

func (s *service) ConfirmAdminRegister(ctx context.Context, email, code string) (*entity.Auth, error) {
	codeFromCache, err := s.cache.GetCode(ctx, email)
	if err != nil {
		return nil, err
	}
	if codeFromCache != code {
		return nil, errors.New("code is invalid")
	}
	passHash, err := s.cache.GetPasswordHash(ctx, email)
	if err != nil {
		return nil, err
	}
	user := entity.NewAdminUser(email, passHash)
	if err := s.usersClient.Create(ctx, user); err != nil {
		return nil, err
	}

	return generateTokens(user, s.accessTTL, s.refreshTTL)
}

func generateTokens(user *entity.User, accessTTL, refreshTTL time.Duration) (*entity.Auth, error) {
	ses := session.NewSession(user.ID, user.Email, user.AccountRole, nil)
	accessToken, err := coreJwt.GenerateJwt(ses, accessTTL, false)
	if err != nil {
		return nil, err
	}
	refreshToken, err := coreJwt.GenerateJwt(ses, refreshTTL, true)
	if err != nil {
		return nil, err
	}
	return &entity.Auth{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func validatePassword(password string) error {
	if len(password) < passwordMinLength {
		return errors.New("password is too short")
	}
	var (
		hasUppercase, hasLowercase, hasDigit, hasSpecial bool
	)
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUppercase = true
		case unicode.IsLower(r):
			hasLowercase = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}
	if !hasUppercase || !hasLowercase || !hasDigit || !hasSpecial {
		return errors.New("password must contain at least one uppercase, lowercase, digit and special character")
	}
	return nil
}

// nolint: govet
func generateCode() string {
	return "1234"
	nums := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	rand.Shuffle(len(nums), func(i, j int) {
		nums[i], nums[j] = nums[j], nums[i]
	})
	return strings.Join(nums[:codeLength], "")
}
