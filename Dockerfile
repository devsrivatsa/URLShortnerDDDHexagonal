# Build stage
FROM golang:1.24.4-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY .env .

# Expose port 8000
EXPOSE 8000:8000

# Command to run the application
CMD ["./main"]
