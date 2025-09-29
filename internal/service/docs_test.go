package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mibrgmv/document-service/internal/domain"
	"github.com/mibrgmv/document-service/internal/service"
	"github.com/mibrgmv/document-service/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDocumentService_UploadDocument_Success_File(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockCacheRepo := new(mocks.MockCacheRepository)
	docService := service.NewDocumentService(mockDocRepo, mockCacheRepo)

	meta := &domain.DocumentMeta{
		Name:   "test.txt",
		File:   true,
		Public: false,
		Mime:   "text/plain",
		Grant:  []string{},
	}
	data := []byte("test file content")
	owner := "testuser"

	mockDocRepo.On("CreateDocument", mock.Anything, mock.MatchedBy(func(doc *domain.Document) bool {
		return doc.Name == "test.txt" &&
			doc.File == true &&
			doc.Owner == "testuser" &&
			string(doc.Data) == "test file content"
	})).Return(nil)

	mockCacheRepo.On("DeletePattern", mock.Anything, "docs:*testuser*").Return(nil)

	doc, err := docService.UploadDocument(context.Background(), meta, data, "", owner)

	assert.NoError(t, err)
	assert.NotNil(t, doc)
	assert.Equal(t, "test.txt", doc.Name)
	assert.True(t, doc.File)
	assert.Equal(t, "testuser", doc.Owner)
	mockDocRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
}

func TestDocumentService_UploadDocument_Success_JSON(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockCacheRepo := new(mocks.MockCacheRepository)
	docService := service.NewDocumentService(mockDocRepo, mockCacheRepo)

	meta := &domain.DocumentMeta{
		Name:   "data.json",
		File:   false,
		Public: true,
		Mime:   "application/json",
		Grant:  []string{"user1", "user2"},
	}
	jsonData := `{"key": "value"}`
	owner := "testuser"

	mockDocRepo.On("CreateDocument", mock.Anything, mock.MatchedBy(func(doc *domain.Document) bool {
		return doc.Name == "data.json" &&
			doc.File == false &&
			doc.Public == true &&
			doc.JSON == `{"key": "value"}` &&
			len(doc.Grant) == 2
	})).Return(nil)

	mockCacheRepo.On("DeletePattern", mock.Anything, "docs:*testuser*").Return(nil)

	doc, err := docService.UploadDocument(context.Background(), meta, nil, jsonData, owner)

	assert.NoError(t, err)
	assert.NotNil(t, doc)
	assert.Equal(t, "data.json", doc.Name)
	assert.False(t, doc.File)
	assert.True(t, doc.Public)
	assert.Equal(t, `{"key": "value"}`, doc.JSON)
	mockDocRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
}

func TestDocumentService_GetDocuments_Success_FromCache(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockCacheRepo := new(mocks.MockCacheRepository)
	docService := service.NewDocumentService(mockDocRepo, mockCacheRepo)

	expectedDocs := []domain.Document{
		{
			ID:      "1",
			Name:    "doc1.txt",
			Mime:    "text/plain",
			File:    true,
			Public:  false,
			Created: time.Now(),
			Grant:   []string{},
		},
	}

	mockCacheRepo.On("GetDocuments", mock.Anything, "docs:testuser::").Return(expectedDocs, nil)

	docs, err := docService.GetDocuments(context.Background(), "testuser", "", "", 100)

	assert.NoError(t, err)
	assert.Equal(t, expectedDocs, docs)
	mockCacheRepo.AssertExpectations(t)
	mockDocRepo.AssertNotCalled(t, "GetUserDocuments")
}

func TestDocumentService_GetDocuments_Success_FromDatabase(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockCacheRepo := new(mocks.MockCacheRepository)
	docService := service.NewDocumentService(mockDocRepo, mockCacheRepo)

	expectedDocs := []domain.Document{
		{
			ID:      "1",
			Name:    "doc1.txt",
			Mime:    "text/plain",
			File:    true,
			Public:  false,
			Created: time.Now(),
			Grant:   []string{},
		},
	}

	mockCacheRepo.On("GetDocuments", mock.Anything, "docs:testuser::").Return(nil, errors.New("cache miss"))
	mockDocRepo.On("GetUserDocuments", mock.Anything, "testuser", 100).Return(expectedDocs, nil)
	mockCacheRepo.On("SetDocuments", mock.Anything, "docs:testuser::", expectedDocs, 5*time.Minute).Return(nil)

	docs, err := docService.GetDocuments(context.Background(), "testuser", "", "", 100)

	assert.NoError(t, err)
	assert.Equal(t, expectedDocs, docs)
	mockDocRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
}

func TestDocumentService_GetDocuments_WithFilter(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockCacheRepo := new(mocks.MockCacheRepository)
	docService := service.NewDocumentService(mockDocRepo, mockCacheRepo)

	allDocs := []domain.Document{
		{
			ID:      "1",
			Name:    "doc1.txt",
			Mime:    "text/plain",
			File:    true,
			Public:  false,
			Created: time.Now(),
			Grant:   []string{},
		},
		{
			ID:      "2",
			Name:    "image.jpg",
			Mime:    "image/jpeg",
			File:    true,
			Public:  true,
			Created: time.Now(),
			Grant:   []string{},
		},
	}

	mockCacheRepo.On("GetDocuments", mock.Anything, "docs:testuser:mime:image/jpeg").Return(nil, errors.New("cache miss"))
	mockDocRepo.On("GetUserDocuments", mock.Anything, "testuser", 100).Return(allDocs, nil)
	mockCacheRepo.On("SetDocuments", mock.Anything, "docs:testuser:mime:image/jpeg", mock.Anything, 5*time.Minute).Return(nil)

	docs, err := docService.GetDocuments(context.Background(), "testuser", "mime", "image/jpeg", 100)

	assert.NoError(t, err)
	assert.Len(t, docs, 1)
	assert.Equal(t, "image.jpg", docs[0].Name)
	assert.Equal(t, "image/jpeg", docs[0].Mime)
	mockDocRepo.AssertExpectations(t)
}

func TestDocumentService_GetDocument_Success_FromCache(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockCacheRepo := new(mocks.MockCacheRepository)
	docService := service.NewDocumentService(mockDocRepo, mockCacheRepo)

	expectedDoc := &domain.Document{
		ID:      "123",
		Name:    "test.txt",
		Mime:    "text/plain",
		File:    true,
		Public:  false,
		Created: time.Now(),
		Grant:   []string{},
		Owner:   "owner1",
		Data:    []byte("content"),
	}

	mockCacheRepo.On("GetDocument", mock.Anything, "doc:123:user123").Return(expectedDoc, nil)

	doc, err := docService.GetDocument(context.Background(), "123", "user123", "testuser")

	assert.NoError(t, err)
	assert.Equal(t, expectedDoc, doc)
	mockCacheRepo.AssertExpectations(t)
	mockDocRepo.AssertNotCalled(t, "GetDocumentByID")
}

func TestDocumentService_GetDocument_Success_FromDatabase(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockCacheRepo := new(mocks.MockCacheRepository)
	docService := service.NewDocumentService(mockDocRepo, mockCacheRepo)

	expectedDoc := &domain.Document{
		ID:      "123",
		Name:    "test.txt",
		Mime:    "text/plain",
		File:    true,
		Public:  true,
		Created: time.Now(),
		Grant:   []string{},
		Owner:   "owner1",
		Data:    []byte("content"),
	}

	mockCacheRepo.On("GetDocument", mock.Anything, "doc:123:user123").Return(nil, errors.New("cache miss"))
	mockDocRepo.On("GetDocumentByID", mock.Anything, "123").Return(expectedDoc, nil)
	mockCacheRepo.On("SetDocument", mock.Anything, "doc:123:user123", expectedDoc, 10*time.Minute).Return(nil)

	doc, err := docService.GetDocument(context.Background(), "123", "user123", "testuser")

	assert.NoError(t, err)
	assert.Equal(t, expectedDoc, doc)
	mockDocRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
}

func TestDocumentService_GetDocument_AccessDenied(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockCacheRepo := new(mocks.MockCacheRepository)
	docService := service.NewDocumentService(mockDocRepo, mockCacheRepo)

	docFromDB := &domain.Document{
		ID:      "123",
		Name:    "private.txt",
		Mime:    "text/plain",
		File:    true,
		Public:  false,
		Created: time.Now(),
		Grant:   []string{},
		Owner:   "otheruser",
		Data:    []byte("content"),
	}

	mockCacheRepo.On("GetDocument", mock.Anything, "doc:123:user123").Return(nil, errors.New("cache miss"))
	mockDocRepo.On("GetDocumentByID", mock.Anything, "123").Return(docFromDB, nil)

	doc, err := docService.GetDocument(context.Background(), "123", "user123", "testuser")

	assert.Error(t, err)
	assert.Nil(t, doc)
	assert.Equal(t, "access denied", err.Error())
	mockDocRepo.AssertExpectations(t)
	mockCacheRepo.AssertNotCalled(t, "SetDocument")
}

func TestDocumentService_GetDocument_AccessGranted_ByGrant(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockCacheRepo := new(mocks.MockCacheRepository)
	docService := service.NewDocumentService(mockDocRepo, mockCacheRepo)

	docFromDB := &domain.Document{
		ID:      "123",
		Name:    "shared.txt",
		Mime:    "text/plain",
		File:    true,
		Public:  false,
		Created: time.Now(),
		Grant:   []string{"testuser"},
		Owner:   "otheruser",
		Data:    []byte("content"),
	}

	mockCacheRepo.On("GetDocument", mock.Anything, "doc:123:user123").Return(nil, errors.New("cache miss"))
	mockDocRepo.On("GetDocumentByID", mock.Anything, "123").Return(docFromDB, nil)
	mockCacheRepo.On("SetDocument", mock.Anything, "doc:123:user123", docFromDB, 10*time.Minute).Return(nil)

	doc, err := docService.GetDocument(context.Background(), "123", "user123", "testuser")

	assert.NoError(t, err)
	assert.Equal(t, docFromDB, doc)
	mockDocRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
}

func TestDocumentService_GetDocument_AccessGranted_ByOwner(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockCacheRepo := new(mocks.MockCacheRepository)
	docService := service.NewDocumentService(mockDocRepo, mockCacheRepo)

	docFromDB := &domain.Document{
		ID:      "123",
		Name:    "myfile.txt",
		Mime:    "text/plain",
		File:    true,
		Public:  false,
		Created: time.Now(),
		Grant:   []string{},
		Owner:   "user123",
		Data:    []byte("content"),
	}

	mockCacheRepo.On("GetDocument", mock.Anything, "doc:123:user123").Return(nil, errors.New("cache miss"))
	mockDocRepo.On("GetDocumentByID", mock.Anything, "123").Return(docFromDB, nil)
	mockCacheRepo.On("SetDocument", mock.Anything, "doc:123:user123", docFromDB, 10*time.Minute).Return(nil)

	doc, err := docService.GetDocument(context.Background(), "123", "user123", "testuser")

	assert.NoError(t, err)
	assert.Equal(t, docFromDB, doc)
	mockDocRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
}

func TestDocumentService_DeleteDocument_Success(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockCacheRepo := new(mocks.MockCacheRepository)
	docService := service.NewDocumentService(mockDocRepo, mockCacheRepo)

	mockDocRepo.On("DeleteDocument", mock.Anything, "123", "testuser").Return(nil)
	mockCacheRepo.On("DeletePattern", mock.Anything, "doc:123*").Return(nil)
	mockCacheRepo.On("DeletePattern", mock.Anything, "docs:*testuser*").Return(nil)

	err := docService.DeleteDocument(context.Background(), "123", "testuser")

	assert.NoError(t, err)
	mockDocRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
}

func TestDocumentService_DeleteDocument_Error(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockCacheRepo := new(mocks.MockCacheRepository)
	docService := service.NewDocumentService(mockDocRepo, mockCacheRepo)

	mockDocRepo.On("DeleteDocument", mock.Anything, "123", "testuser").Return(errors.New("database error"))

	err := docService.DeleteDocument(context.Background(), "123", "testuser")

	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	mockDocRepo.AssertExpectations(t)
	mockCacheRepo.AssertNotCalled(t, "DeletePattern")
}

func TestDocumentService_FilterDocuments(t *testing.T) {
	mockDocRepo := new(mocks.MockDocumentRepository)
	mockCacheRepo := new(mocks.MockCacheRepository)
	docService := service.NewDocumentService(mockDocRepo, mockCacheRepo)

	docs := []domain.Document{
		{Name: "doc1.txt", Mime: "text/plain", Public: false},
		{Name: "doc2.pdf", Mime: "application/pdf", Public: true},
		{Name: "image.jpg", Mime: "image/jpeg", Public: true},
	}

	// Filter by name
	filtered := docService.FilterDocuments(docs, "name", "doc1.txt")
	assert.Len(t, filtered, 1)
	assert.Equal(t, "doc1.txt", filtered[0].Name)

	// Filter by mime
	filtered = docService.FilterDocuments(docs, "mime", "image/jpeg")
	assert.Len(t, filtered, 1)
	assert.Equal(t, "image.jpg", filtered[0].Name)

	// Filter by public=true
	filtered = docService.FilterDocuments(docs, "public", "true")
	assert.Len(t, filtered, 2)

	// Filter by public=false
	filtered = docService.FilterDocuments(docs, "public", "false")
	assert.Len(t, filtered, 1)
	assert.Equal(t, "doc1.txt", filtered[0].Name)

	// Unknown filter key - returns empty
	filtered = docService.FilterDocuments(docs, "unknown", "value")
	assert.Len(t, filtered, 0)
}
