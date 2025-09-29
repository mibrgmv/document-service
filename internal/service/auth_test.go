package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/mibrgmv/document-service/internal/service"
	"github.com/mibrgmv/document-service/internal/service/mocks"
	"github.com/mibrgmv/document-service/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_Register_Success(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockCacheRepo := new(mocks.MockCacheRepository)
	jwtManager := jwt.NewManager("test-secret", 24*time.Hour)

	authService := service.NewAuthService(mockUserRepo, mockCacheRepo, jwtManager, "admin-token")

	mockUserRepo.On("UserExists", mock.Anything, "testuser").Return(false, nil)
	mockUserRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*domain.User")).Return(nil)

	err := authService.Register(context.Background(),
		"admin-token",
		"testuser",
		"Password123!",
	)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_Register_InvalidAdminToken(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockCacheRepo := new(mocks.MockCacheRepository)
	jwtManager := jwt.NewManager("test-secret", 24*time.Hour)

	authService := service.NewAuthService(mockUserRepo, mockCacheRepo, jwtManager, "admin-token")

	err := authService.Register(context.Background(),
		"wrong-token",
		"testuser",
		"Password123!",
	)

	assert.Error(t, err)
	assert.Equal(t, "invalid admin token", err.Error())
}

func TestAuthService_Register_UserExists(t *testing.T) {
	mockUserRepo := new(mocks.MockUserRepository)
	mockCacheRepo := new(mocks.MockCacheRepository)
	jwtManager := jwt.NewManager("test-secret", 24*time.Hour)

	authService := service.NewAuthService(mockUserRepo, mockCacheRepo, jwtManager, "admin-token")

	mockUserRepo.On("UserExists", mock.Anything, "existinguser").Return(true, nil)

	err := authService.Register(context.Background(),
		"admin-token",
		"existinguser",
		"Password123!",
	)

	assert.Error(t, err)
	assert.Equal(t, "user already exists", err.Error())
	mockUserRepo.AssertExpectations(t)
}
