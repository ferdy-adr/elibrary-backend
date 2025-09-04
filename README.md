# eLibrary Backend

Backend API untuk aplikasi perpustakaan mobile menggunakan Go dengan framework Gin.

## Features

- ✅ **Authentication & Authorization** - JWT-based authentication (no expiry)
- ✅ **Book Management** - CRUD operations untuk buku
- ✅ **File Upload** - Upload gambar cover buku
- ✅ **Search & Filter** - Pencarian dan filter buku berdasarkan berbagai kriteria
- ✅ **Pagination** - Pagination untuk list buku
- ✅ **Database Migration** - Database schema management
- ✅ **Docker Support** - Containerized MySQL database

## Tech Stack

- **Language**: Go 1.23.5
- **Framework**: Gin
- **Database**: MySQL 8.0
- **Authentication**: JWT
- **File Upload**: Multipart form-data
- **Migration**: golang-migrate
- **Containerization**: Docker & Docker Compose

## Project Structure

```
.
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── configs/               # Configuration management
│   │   ├── config.go
│   │   ├── config.yaml
│   │   └── types.go
│   ├── handlers/              # HTTP handlers
│   │   ├── auth/
│   │   └── books/
│   ├── middleware/            # Middleware
│   │   ├── auth.go
│   │   └── cors.go
│   ├── model/                 # Data models
│   │   ├── book.go
│   │   ├── response.go
│   │   └── user.go
│   ├── repository/            # Data access layer
│   │   ├── books/
│   │   └── users/
│   └── service/               # Business logic layer
│       ├── auth/
│       └── books/
├── pkg/
│   └── internalsql/           # Database connection
│       └── sql.go
├── scripts/
│   └── migrations/            # Database migrations
│       ├── 000001_create_users_table.up.sql
│       ├── 000001_create_users_table.down.sql
│       ├── 000002_create_books_table.up.sql
│       └── 000002_create_books_table.down.sql
├── public/
│   └── images/                # Uploaded cover images
├── docker-compose.yml         # Docker configuration
├── Makefile                   # Build and migration commands
├── go.mod                     # Go modules
└── API_DOCUMENTATION.md       # API documentation
```

## Setup & Installation

### Prerequisites

- Go 1.23.5 or later
- Docker & Docker Compose
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI tool

### 1. Clone Repository

```bash
git clone <repository-url>
cd elibrary-backend
```

### 2. Install Dependencies

```bash
make install
```

### 3. Setup Database

```bash
# Start MySQL container
make docker-up

# Run database migrations
make migrate-up
```

### 4. Configure Application

Update `internal/configs/config.yaml` if needed:

```yaml
service:
  port: ":8080"

database:
  dataSourceName: "root:secretPassword@tcp(localhost:3306)/elibrary"

jwt:
  secretKey: "your-very-secret-key-for-elibrary-app"

upload:
  path: "./public/images"
```

### 5. Run Application

```bash
# Development mode (includes docker-up, migrate-up, and run)
make dev

# Or run individually
make run
```

## Available Commands

### Development
```bash
make dev          # Start development environment
make run          # Run application
make build        # Build binary
make test         # Run tests
make install      # Install dependencies
```

### Docker
```bash
make docker-up    # Start Docker containers
make docker-down  # Stop Docker containers
make docker-logs  # View Docker logs
```

### Database Migrations
```bash
make migrate-up                    # Run all migrations
make migrate-down                  # Rollback all migrations
make migrate-create name=new_table # Create new migration
make migrate-force version=1       # Force migration version
```

## API Usage

### 1. Register User

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com", 
    "password": "password123",
    "full_name": "John Doe"
  }'
```

### 2. Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "password123"
  }'
```

### 3. Get Books (Public)

```bash
curl "http://localhost:8080/api/books?page=1&limit=10&search=harry"
```

### 4. Create Book (Protected)

```bash
curl -X POST http://localhost:8080/api/books \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "title=Harry Potter" \
  -F "isbn=978-0747532699" \
  -F "year=1997" \
  -F "publisher=Bloomsbury" \
  -F "author=J.K. Rowling" \
  -F "synopsis=A young wizard story..." \
  -F "cover_image=@/path/to/cover.jpg"
```

### 5. Update Book (Protected)

```bash
curl -X PATCH http://localhost:8080/api/books/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "title=Updated Title" \
  -F "cover_image=@/path/to/new_cover.jpg"
```

### 6. Delete Book (Protected)

```bash
curl -X DELETE http://localhost:8080/api/books/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## API Documentation

Lihat [API_DOCUMENTATION.md](./API_DOCUMENTATION.md) untuk dokumentasi lengkap API endpoints.

## Environment Variables

Aplikasi ini menggunakan file konfigurasi YAML, namun Anda juga bisa menggunakan environment variables:

```bash
export SERVICE_PORT=":8080"
export DATABASE_DATASOURCENAME="root:secretPassword@tcp(localhost:3306)/elibrary"
export JWT_SECRETKEY="your-secret-key"
export UPLOAD_PATH="./public/images"
```

## Database Schema

### Users Table
- `id` (INT, Primary Key, Auto Increment)
- `username` (VARCHAR, Unique, Not Null)
- `email` (VARCHAR, Unique, Not Null) 
- `password` (VARCHAR, Not Null) - Hashed with bcrypt
- `full_name` (VARCHAR, Not Null)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

### Books Table
- `id` (INT, Primary Key, Auto Increment)
- `title` (VARCHAR, Not Null)
- `isbn` (VARCHAR, Unique, Not Null)
- `year` (INT, Not Null)
- `publisher` (VARCHAR, Not Null)
- `author` (VARCHAR, Not Null)
- `cover_image` (VARCHAR, Nullable)
- `synopsis` (TEXT, Nullable)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

## Security Features

- **Password Hashing**: Menggunakan bcrypt untuk hash password
- **JWT Authentication**: Token untuk autentikasi API
- **CORS**: Cross-Origin Resource Sharing support
- **File Validation**: Validasi tipe file untuk upload gambar
- **SQL Injection Protection**: Menggunakan prepared statements

## File Upload

- **Supported Formats**: JPG, JPEG, PNG
- **Storage Location**: `./public/images/`
- **File Naming**: `cover_{timestamp}.{extension}`
- **Access URL**: `http://localhost:8080/images/{filename}`

## Troubleshooting

### Database Connection Issues
```bash
# Check if MySQL container is running
docker ps

# Check logs
make docker-logs

# Recreate containers
make docker-down
make docker-up
```

### Migration Issues
```bash
# Check migration status
migrate -database "mysql://root:secretPassword@tcp(localhost:3306)/elibrary" -path scripts/migrations version

# Force specific version
make migrate-force version=1
```

### Permission Issues
```bash
# Make sure public/images directory exists and is writable
mkdir -p public/images
chmod 755 public/images
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.
