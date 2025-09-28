package repository

import (
	"context"

	"github.com/mibrgmv/document-service/internal/domain"
)

type DocumentRepository interface {
	CreateDocument(ctx context.Context, doc *domain.Document) error
	GetDocumentByID(ctx context.Context, id string) (*domain.Document, error)
	GetUserDocuments(ctx context.Context, login string, limit int) ([]domain.Document, error)
	DeleteDocument(ctx context.Context, id, owner string) error
	DocumentExists(ctx context.Context, id string) (bool, error)
}
