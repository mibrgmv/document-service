package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mibrgmv/document-service/internal/service"
)

type DocumentHandler struct {
	docService service.DocumentService
}

func NewDocumentHandler(docService service.DocumentService) *DocumentHandler {
	return &DocumentHandler{docService: docService}
}

// UploadDocument godoc
// @Summary Upload document
// @Description Upload a new document (file or JSON)
// @Tags documents
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param meta formData string true "Document metadata JSON"
// @Param file formData file false "Document file"
// @Param json formData string false "JSON data (if not file)"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Router /docs [post]
func (h *DocumentHandler) UploadDocument(c *gin.Context) {
	login := c.MustGet("login").(string)

	metaStr := c.PostForm("meta")
	var meta DocumentMeta
	if err := json.Unmarshal([]byte(metaStr), &meta); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Error: &Error{Code: 400, Text: "invalid meta data"},
		})
		return
	}

	file, err := c.FormFile("file")
	var fileData []byte
	if err == nil {
		file, err := file.Open()
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Error: &Error{Code: 400, Text: "invalid file"},
			})
			return
		}
		defer file.Close()

		fileData, err = io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Error: &Error{Code: 400, Text: "invalid file"},
			})
			return
		}
	}

	jsonData := c.PostForm("json")

	doc, err := h.docService.UploadDocument(c.Request.Context(), meta.ToDomain(), fileData, jsonData, login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Error: &Error{Code: 500, Text: err.Error()},
		})
		return
	}

	responseData := gin.H{}
	if doc.File {
		responseData["file"] = doc.Name
	} else {
		responseData["json"] = doc.JSON
	}

	c.JSON(http.StatusOK, Response{
		Data: responseData,
	})
}

// GetDocuments godoc
// @Summary Get documents list
// @Description Get list of documents with optional filtering
// @Tags documents
// @Security BearerAuth
// @Produce json
// @Param login query string false "User login to filter (default: current user)"
// @Param key query string false "Filter key (name, mime, public)"
// @Param value query string false "Filter value"
// @Param limit query integer false "Limit number of documents"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Failure 500 {object} Response
// @Router /docs [get]
func (h *DocumentHandler) GetDocuments(c *gin.Context) {
	login := c.MustGet("login").(string)
	targetLogin := c.Query("login")
	if targetLogin == "" {
		targetLogin = login
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	filterKey := c.Query("key")
	filterValue := c.Query("value")

	docs, err := h.docService.GetDocuments(c.Request.Context(), targetLogin, filterKey, filterValue, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Error: &Error{Code: 500, Text: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Data: gin.H{"docs": docs},
	})
}

// GetDocumentsHead godoc
// @Summary HEAD documents list
// @Description HEAD request for documents list
// @Tags documents
// @Security BearerAuth
// @Router /docs [head]
func (h *DocumentHandler) GetDocumentsHead(c *gin.Context) {
	c.Status(http.StatusOK)
}

// GetDocument godoc
// @Summary Get document
// @Description Get document by ID. Returns file or JSON based on document type
// @Tags documents
// @Security BearerAuth
// @Produce json,application/octet-stream
// @Param id path string true "Document ID"
// @Success 200 {object} Response
// @Failure 403 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /docs/{id} [get]
func (h *DocumentHandler) GetDocument(c *gin.Context) {
	login := c.MustGet("login").(string)
	id := c.Param("id")

	doc, err := h.docService.GetDocument(c.Request.Context(), id, login)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "access denied" {
			status = http.StatusForbidden
		}
		c.JSON(status, Response{
			Error: &Error{Code: status, Text: err.Error()},
		})
		return
	}

	if doc.File {
		c.Header("Content-Type", doc.Mime)
		c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", doc.Name))
		c.Header("Content-Length", strconv.Itoa(len(doc.Data)))
		c.Data(http.StatusOK, doc.Mime, doc.Data)
		return
	}

	c.JSON(http.StatusOK, Response{
		Data: doc.JSON,
	})
}

// GetDocumentHead godoc
// @Summary HEAD document
// @Description HEAD request for document
// @Tags documents
// @Security BearerAuth
// @Param id path string true "Document ID"
// @Router /docs/{id} [head]
func (h *DocumentHandler) GetDocumentHead(c *gin.Context) {
	c.Status(http.StatusOK)
}

// DeleteDocument godoc
// @Summary Delete document
// @Description Delete document by ID
// @Tags documents
// @Security BearerAuth
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Failure 403 {object} Response
// @Failure 500 {object} Response
// @Router /docs/{id} [delete]
func (h *DocumentHandler) DeleteDocument(c *gin.Context) {
	login := c.MustGet("login").(string)
	id := c.Param("id")

	if err := h.docService.DeleteDocument(c.Request.Context(), id, login); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Error: &Error{Code: 500, Text: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Response: gin.H{id: true},
	})
}
