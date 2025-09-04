package books

import (
	"net/http"
	"strconv"

	"github.com/ferdy-adr/elibrary-backend/internal/middleware"
	"github.com/ferdy-adr/elibrary-backend/internal/model"
	bookService "github.com/ferdy-adr/elibrary-backend/internal/service/books"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	bookService *bookService.Service
}

func NewHandler(bookService *bookService.Service) *Handler {
	return &Handler{
		bookService: bookService,
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	// Public routes (for reading books)
	public := r.Group("/api/books")
	{
		public.GET("", h.GetBooks)
		public.GET("/:id", h.GetBookByID)
	}

	// Protected routes (for managing books)
	protected := r.Group("/api/books")
	protected.Use(middleware.JWTMiddleware())
	{
		protected.POST("", h.CreateBook)
		protected.PATCH("/:id", h.UpdateBook)
		protected.DELETE("/:id", h.DeleteBook)
	}

	// Static files for images
	r.Static("/images", "./public/images")
}

func (h *Handler) GetBooks(c *gin.Context) {
	var params model.BookQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
		return
	}

	response, err := h.bookService.GetBooks(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Message: "Failed to get books",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "Books retrieved successfully",
		Data:    response,
	})
}

func (h *Handler) GetBookByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Invalid book ID",
			Error:   "Book ID must be a number",
		})
		return
	}

	book, err := h.bookService.GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIResponse{
			Success: false,
			Message: "Book not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "Book retrieved successfully",
		Data:    book,
	})
}

func (h *Handler) CreateBook(c *gin.Context) {
	var req model.CreateBookRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	// Handle cover image upload
	coverFile, _ := c.FormFile("cover_image")

	book, err := h.bookService.CreateBook(req, coverFile)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "ISBN already exists" {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, model.APIResponse{
			Success: false,
			Message: "Failed to create book",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, model.APIResponse{
		Success: true,
		Message: "Book created successfully",
		Data:    book,
	})
}

func (h *Handler) UpdateBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Invalid book ID",
			Error:   "Book ID must be a number",
		})
		return
	}

	var req model.UpdateBookRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	// Handle cover image upload
	coverFile, _ := c.FormFile("cover_image")

	book, err := h.bookService.UpdateBook(id, req, coverFile)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "book not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "ISBN already exists" {
			statusCode = http.StatusConflict
		}

		c.JSON(statusCode, model.APIResponse{
			Success: false,
			Message: "Failed to update book",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "Book updated successfully",
		Data:    book,
	})
}

func (h *Handler) DeleteBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Message: "Invalid book ID",
			Error:   "Book ID must be a number",
		})
		return
	}

	err = h.bookService.DeleteBook(id)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "book not found" {
			statusCode = http.StatusNotFound
		}

		c.JSON(statusCode, model.APIResponse{
			Success: false,
			Message: "Failed to delete book",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "Book deleted successfully",
	})
}
