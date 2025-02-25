# Use the official Golang image as a build stage
FROM golang:1.23 AS builder

# Set environment variables for building
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to leverage Go module caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

#RUN go install github.com/air-verse/air@latest

# Build the Go application
RUN go build -o main ./cmd/producer/main.go


# Use a minimal runtime image for the final container
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose port 8002
EXPOSE 8002

# Run the application
CMD ["./main"]
