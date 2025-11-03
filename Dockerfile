# Step 1: Build the Go app
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary
RUN go build -o main .

# Step 2: Create a small final image
FROM alpine:latest

WORKDIR /root/

# Copy the compiled binary from builder
COPY --from=builder /app/main .

# Expose port (your Go app port)
EXPOSE 8080

# Run the Go app
CMD ["./main"]
