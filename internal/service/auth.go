package service

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/mibrgmv/document-service/internal/domain"
	"github.com/mibrgmv/document-service/internal/repository"
	"github.com/mibrgmv/document-service/pkg/jwt"
	"github.com/mibrgmv/document-service/pkg/utils"
)

type AuthService interface {
	Register(ctx context.Context, token, login, password string) error
	Authenticate(ctx context.Context, login, password string) (string, error)
	Logout(ctx context.Context, token string) error
}

type authService struct {
	userRepo   repository.UserRepository
	cacheRepo  repository.CacheRepository
	jwtManager *jwt.Manager
	adminToken string
}

func NewAuthService(
	userRepo repository.UserRepository,
	cacheRepo repository.CacheRepository,
	jwtManager *jwt.Manager,
	adminToken string,
) AuthService {
	return &authService{
		userRepo:   userRepo,
		cacheRepo:  cacheRepo,
		jwtManager: jwtManager,
		adminToken: adminToken,
	}
}

func (s *authService) Register(ctx context.Context, token, login, password string) error {
	if token != s.adminToken {
		return errors.New("invalid admin token")
	}

	if !isValidLogin(login) {
		return errors.New("login must be at least 4 characters long and contain only letters and numbers")
	}

	if !isValidPassword(password) {
		return errors.New("password must be at least 4 characters long, contain uppercase and lowercase letter, digit and special character")
	}

	exists, err := s.userRepo.UserExists(ctx, login)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	user := &domain.User{
		ID:       utils.GenerateID(),
		Login:    login,
		Password: hashedPassword,
		Created:  time.Now(),
	}

	return s.userRepo.CreateUser(ctx, user)
}

func (s *authService) Authenticate(ctx context.Context, login, password string) (string, error) {
	user, err := s.userRepo.GetUserByLogin(ctx, login)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	token, err := s.jwtManager.GenerateToken(user.ID, user.Login)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) Logout(ctx context.Context, token string) error {
	return s.cacheRepo.DeletePattern(ctx, "*"+token+"*")
}

func isValidLogin(login string) bool {
	if len(login) < 4 {
		return false
	}
	matched, _ := regexp.MatchString("^[a-zA-Z0-9]+$", login)
	return matched
}

func isValidPassword(password string) bool {
	if len(password) < 4 {
		return false
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)

	return hasUpper && hasLower && hasDigit && hasSpecial
}
