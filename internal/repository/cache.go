package repository

import (
	"context"
	"time"

	"github.com/mibrgmv/document-service/internal/domain"
)

type CacheRepository interface {
	SetDocuments(ctx context.Context, key string, docs []domain.Document, expiration time.Duration) error
	GetDocuments(ctx context.Context, key string) ([]domain.Document, error)
	SetDocument(ctx context.Context, key string, doc *domain.Document, expiration time.Duration) error
	GetDocument(ctx context.Context, key string) (*domain.Document, error)
	Delete(ctx context.Context, key string) error
	DeletePattern(ctx context.Context, pattern string) error
}
