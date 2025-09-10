# Build stage
FROM golang:1.20-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/schema.sql .
COPY --from=builder /app/migrations/002_seed_data.sql migrations/

# Create directory for database
RUN mkdir -p /data

# Set environment variables
ENV DATABASE_URL=/data/auth.db
ENV GIN_MODE=release

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
