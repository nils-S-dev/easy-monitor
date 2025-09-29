# Stage 1: Build the Go binary
FROM golang:1.25.1 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o easy-monitor ./cmd

# Stage 2: Create a minimal final image
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/easy-monitor .

# Expose port (if your app listens on one)
EXPOSE 8080

# Command to run the executable
CMD ["./easy-monitor"]