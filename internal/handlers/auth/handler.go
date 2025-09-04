package auth

import (
	"net/http"

	"github.com/ferdy-adr/elibrary-backend/internal/model"
	authService "github.com/ferdy-adr/elibrary-backend/internal/service/auth"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	authService *authService.Service
}

func NewHandler(authService *authService.Service) *Handler {
	return &Handler{
		authService: authService,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	auth := r.Group("/api/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}
}

func (h *Handler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "username already exists" {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, model.APIResponse{
			Success: false,
			Message: "Registration failed",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, model.APIResponse{
		Success: true,
		Message: "User registered successfully",
		Data:    user,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	loginResponse, err := h.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.APIResponse{
			Success: false,
			Message: "Login failed",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "Login successful",
		Data:    loginResponse,
	})
}
