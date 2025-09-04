# Build stage
FROM golang:1.23.5-alpine AS builder

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates
RUN apk --no-cache add ca-certificates

# Create app directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy configuration files
COPY --from=builder /app/internal/configs/ ./internal/configs/

# Copy scripts for migrations
COPY --from=builder /app/scripts/ ./scripts/

# Create upload directory
RUN mkdir -p ./public/images

# Expose port
EXPOSE 8080

# Command to run
CMD ["./main"]
