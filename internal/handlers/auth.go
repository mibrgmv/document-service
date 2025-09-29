package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mibrgmv/document-service/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary Register new user
// @Description Register a new user (requires admin token)
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration data"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Router /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Error: &Error{Code: 400, Text: "invalid request"},
		})
		return
	}

	if err := h.authService.Register(c.Request.Context(), req.Token, req.Login, req.Pswd); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Error: &Error{Code: 400, Text: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Response: gin.H{"login": req.Login},
	})
}

// Auth godoc
// @Summary User authentication
// @Description Authenticate user and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body AuthRequest true "Auth credentials"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 401 {object} Response
// @Router /auth [post]
func (h *AuthHandler) Auth(c *gin.Context) {
	var req AuthRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Error: &Error{Code: 400, Text: "invalid request"},
		})
		return
	}

	token, err := h.authService.Authenticate(c.Request.Context(), req.Login, req.Pswd)
	if err != nil {
		c.JSON(http.StatusUnauthorized, Response{
			Error: &Error{Code: 401, Text: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Response: gin.H{"token": token},
	})
}

// Logout godoc
// @Summary User logout
// @Description Logout user and invalidate token
// @Tags auth
// @Produce json
// @Param token path string true "JWT Token"
// @Success 200 {object} Response
// @Failure 500 {object} Response
// @Router /auth/{token} [delete]
func (h *AuthHandler) Logout(c *gin.Context) {
	token := c.Param("token")
	if err := h.authService.Logout(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Error: &Error{Code: 500, Text: "internal error"},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Response: gin.H{token: true},
	})
}
