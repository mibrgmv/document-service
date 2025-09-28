package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mibrgmv/document-service/internal/domain"
	"github.com/mibrgmv/document-service/internal/repository"
)

type cacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) repository.CacheRepository {
	return &cacheRepository{client: client}
}

func (r *cacheRepository) SetDocuments(ctx context.Context, key string, docs []domain.Document, expiration time.Duration) error {
	data, err := json.Marshal(docs)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, expiration).Err()
}

func (r *cacheRepository) GetDocuments(ctx context.Context, key string) ([]domain.Document, error) {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var docs []domain.Document
	err = json.Unmarshal(data, &docs)
	return docs, err
}

func (r *cacheRepository) SetDocument(ctx context.Context, key string, doc *domain.Document, expiration time.Duration) error {
	data, err := json.Marshal(doc)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, expiration).Err()
}

func (r *cacheRepository) GetDocument(ctx context.Context, key string) (*domain.Document, error) {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var doc domain.Document
	err = json.Unmarshal(data, &doc)
	return &doc, err
}

func (r *cacheRepository) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *cacheRepository) DeletePattern(ctx context.Context, pattern string) error {
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return r.client.Del(ctx, keys...).Err()
	}
	return nil
}
