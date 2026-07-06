# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /tigercasino ./cmd/server

# Production stage
FROM alpine:latest

WORKDIR /app

# Install certificates
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /tigercasino .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./tigercasino"]
