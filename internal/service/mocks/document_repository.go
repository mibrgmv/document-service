package mocks

import (
	"context"

	"github.com/mibrgmv/document-service/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockDocumentRepository struct {
	mock.Mock
}

func (m *MockDocumentRepository) CreateDocument(ctx context.Context, doc *domain.Document) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

func (m *MockDocumentRepository) GetDocumentByID(ctx context.Context, id string) (*domain.Document, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Document), args.Error(1)
}

func (m *MockDocumentRepository) GetUserDocuments(ctx context.Context, login string, limit int) ([]domain.Document, error) {
	args := m.Called(ctx, login, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Document), args.Error(1)
}

func (m *MockDocumentRepository) DeleteDocument(ctx context.Context, id, owner string) error {
	args := m.Called(ctx, id, owner)
	return args.Error(0)
}

func (m *MockDocumentRepository) DocumentExists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}
