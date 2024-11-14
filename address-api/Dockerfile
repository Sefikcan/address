# Stage 1: Build the Go binary
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . /app

# Build the application binary
RUN go build -o main ./cmd/server/main.go

# Stage 2: Create a smaller image for running the application
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy the config directory
COPY --from=builder /app/pkg/config /root/pkg/config

# Expose the port on which the app will run
EXPOSE 3048

# Command to run the binary
CMD ["./main"]
