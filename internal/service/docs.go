package service

import (
	"context"
	"errors"
	"time"

	"github.com/mibrgmv/document-service/internal/domain"
	"github.com/mibrgmv/document-service/internal/repository"
	"github.com/mibrgmv/document-service/pkg/utils"
)

type DocumentService interface {
	UploadDocument(ctx context.Context, meta *domain.DocumentMeta, data []byte, jsonData, owner string) (*domain.Document, error)
	GetDocuments(ctx context.Context, login, filterKey, filterValue string, limit int) ([]domain.Document, error)
	GetDocument(ctx context.Context, id, login string) (*domain.Document, error)
	DeleteDocument(ctx context.Context, id, owner string) error
	FilterDocuments(docs []domain.Document, key, value string) []domain.Document
}

type documentService struct {
	docRepo   repository.DocumentRepository
	cacheRepo repository.CacheRepository
}

func NewDocumentService(
	docRepo repository.DocumentRepository,
	cacheRepo repository.CacheRepository,
) DocumentService {
	return &documentService{
		docRepo:   docRepo,
		cacheRepo: cacheRepo,
	}
}

func (s *documentService) UploadDocument(ctx context.Context, meta *domain.DocumentMeta, data []byte, jsonData, owner string) (*domain.Document, error) {
	doc := &domain.Document{
		ID:      utils.GenerateID(),
		Name:    meta.Name,
		Mime:    meta.Mime,
		File:    meta.File,
		Public:  meta.Public,
		Created: time.Now(),
		Grant:   meta.Grant,
		Owner:   owner,
	}

	if meta.File {
		doc.Data = data
	} else {
		doc.JSON = jsonData
	}

	err := s.docRepo.CreateDocument(ctx, doc)
	if err != nil {
		return nil, err
	}

	s.cacheRepo.DeletePattern(ctx, "docs:*"+owner+"*")
	return doc, nil
}

func (s *documentService) GetDocuments(ctx context.Context, login, filterKey, filterValue string, limit int) ([]domain.Document, error) {
	cacheKey := "docs:" + login + ":" + filterKey + ":" + filterValue

	if cached, err := s.cacheRepo.GetDocuments(ctx, cacheKey); err == nil {
		return cached, nil
	}

	docs, err := s.docRepo.GetUserDocuments(ctx, login, limit)
	if err != nil {
		return nil, err
	}

	if filterKey != "" && filterValue != "" {
		docs = s.FilterDocuments(docs, filterKey, filterValue)
	}

	s.cacheRepo.SetDocuments(ctx, cacheKey, docs, 5*time.Minute)
	return docs, nil
}

func (s *documentService) GetDocument(ctx context.Context, id, login string) (*domain.Document, error) {
	cacheKey := "doc:" + id + ":" + login

	if cached, err := s.cacheRepo.GetDocument(ctx, cacheKey); err == nil {
		return cached, nil
	}

	doc, err := s.docRepo.GetDocumentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if doc.Owner != login && !doc.Public && !contains(doc.Grant, login) {
		return nil, errors.New("access denied")
	}

	s.cacheRepo.SetDocument(ctx, cacheKey, doc, 10*time.Minute)
	return doc, nil
}

func (s *documentService) DeleteDocument(ctx context.Context, id, owner string) error {
	err := s.docRepo.DeleteDocument(ctx, id, owner)
	if err != nil {
		return err
	}

	s.cacheRepo.DeletePattern(ctx, "doc:"+id+"*")
	s.cacheRepo.DeletePattern(ctx, "docs:*"+owner+"*")
	return nil
}

func (s *documentService) FilterDocuments(docs []domain.Document, key, value string) []domain.Document {
	var filtered []domain.Document
	for _, doc := range docs {
		switch key {
		case "name":
			if doc.Name == value {
				filtered = append(filtered, doc)
			}
		case "mime":
			if doc.Mime == value {
				filtered = append(filtered, doc)
			}
		case "public":
			if doc.Public == (value == "true") {
				filtered = append(filtered, doc)
			}
		}
	}
	return filtered
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
