package main

import (
	"log"

	"github.com/ferdy-adr/elibrary-backend/internal/configs"
	authHandler "github.com/ferdy-adr/elibrary-backend/internal/handlers/auth"
	bookHandler "github.com/ferdy-adr/elibrary-backend/internal/handlers/books"
	"github.com/ferdy-adr/elibrary-backend/internal/middleware"
	bookRepo "github.com/ferdy-adr/elibrary-backend/internal/repository/books"
	userRepo "github.com/ferdy-adr/elibrary-backend/internal/repository/users"
	authService "github.com/ferdy-adr/elibrary-backend/internal/service/auth"
	bookService "github.com/ferdy-adr/elibrary-backend/internal/service/books"
	"github.com/ferdy-adr/elibrary-backend/pkg/internalsql"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize config
	err := configs.Init(
		configs.WithConfigFolder([]string{"./internal/configs"}),
		configs.WithConfigFile("config"),
		configs.WithConfigType("yaml"),
	)
	if err != nil {
		log.Fatal("Failed to initialize config:", err)
	}

	cfg := configs.Get()
	log.Println("Config loaded successfully")

	// Initialize database
	db, err := internalsql.Connect(cfg.Database.DataSourceName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connected successfully")

	// Initialize repositories
	userRepository := userRepo.NewRepository(db)
	bookRepository := bookRepo.NewRepository(db)

	// Initialize services
	authSvc := authService.NewService(userRepository)
	bookSvc := bookService.NewService(bookRepository)

	// Initialize handlers
	authHdl := authHandler.NewHandler(authSvc)
	bookHdl := bookHandler.NewHandler(bookSvc)

	// Initialize Gin router
	r := gin.Default()

	// Add middleware
	r.Use(middleware.CORSMiddleware())

	// Register routes
	authHdl.RegisterRoutes(r)
	bookHdl.RegisterRoutes(r)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "eLibrary Backend is running",
		})
	})

	// Start server
	log.Printf("Server starting on port %s", cfg.Service.Port)
	if err := r.Run(cfg.Service.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
