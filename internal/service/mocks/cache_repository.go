package mocks

import (
	"context"
	"time"

	"github.com/mibrgmv/document-service/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockCacheRepository struct {
	mock.Mock
}

func (m *MockCacheRepository) SetDocuments(ctx context.Context, key string, docs []domain.Document, expiration time.Duration) error {
	args := m.Called(ctx, key, docs, expiration)
	return args.Error(0)
}

func (m *MockCacheRepository) GetDocuments(ctx context.Context, key string) ([]domain.Document, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Document), args.Error(1)
}

func (m *MockCacheRepository) SetDocument(ctx context.Context, key string, doc *domain.Document, expiration time.Duration) error {
	args := m.Called(ctx, key, doc, expiration)
	return args.Error(0)
}

func (m *MockCacheRepository) GetDocument(ctx context.Context, key string) (*domain.Document, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Document), args.Error(1)
}

func (m *MockCacheRepository) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCacheRepository) DeletePattern(ctx context.Context, pattern string) error {
	args := m.Called(ctx, pattern)
	return args.Error(0)
}
