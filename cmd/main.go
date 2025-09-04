package main

import (
	"database/sql"
	"log"
	"net/url"
	"os"
	"strings"

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
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// convertMySQLURL converts Railway MySQL URL format to Go driver format
// From: mysql://user:password@host:port/database
// To: user:password@tcp(host:port)/database
func convertMySQLURL(mysqlURL string) (string, error) {
	if !strings.HasPrefix(mysqlURL, "mysql://") {
		// Already in correct format
		return mysqlURL, nil
	}

	u, err := url.Parse(mysqlURL)
	if err != nil {
		return "", err
	}

	user := u.User.Username()
	password, _ := u.User.Password()
	host := u.Host
	database := strings.TrimPrefix(u.Path, "/")

	// Convert to Go MySQL driver format
	dsn := user + ":" + password + "@tcp(" + host + ")/" + database
	return dsn, nil
}

// runMigrations runs database migrations automatically
func runMigrations(db *sql.DB) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./scripts/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}

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

	// Debug info for Railway deployment
	log.Printf("Database DSN: %s", cfg.Database.DataSourceName)
	log.Printf("Environment PORT: %s", os.Getenv("PORT"))
	log.Printf("Environment DATABASE_URL: %s", os.Getenv("DATABASE_URL"))

	// Handle Railway's PORT environment variable
	port := cfg.Service.Port
	if railwayPort := os.Getenv("PORT"); railwayPort != "" {
		port = ":" + railwayPort
	}

	// Handle Railway's DATABASE_URL environment variable
	dataSourceName := cfg.Database.DataSourceName
	if railwayDataSource := os.Getenv("DATABASE_URL"); railwayDataSource != "" {
		convertedDSN, err := convertMySQLURL(railwayDataSource)
		if err != nil {
			log.Fatal("Failed to convert DATABASE_URL:", err)
		}
		dataSourceName = convertedDSN
		log.Printf("Using Railway DATABASE_URL: %s", railwayDataSource)
		log.Printf("Converted DSN: %s", dataSourceName)
	}

	// Set Gin mode for production
	if os.Getenv("GIN_MODE") == "" && os.Getenv("PORT") != "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	db, err := internalsql.Connect(dataSourceName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connected successfully")

	// Run migrations automatically
	if err := runMigrations(db); err != nil {
		log.Printf("Warning: Migration failed: %v", err)
		// Don't fatal here, let the app continue if migrations fail
	}

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
	log.Printf("Server starting on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
