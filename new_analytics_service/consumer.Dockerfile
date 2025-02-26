# Use the official Golang image as a build stage
FROM golang:1.23-alpine AS builder

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


# Build the Go application
RUN go build -o main ./cmd/consumer/main.go


# Use a minimal runtime image for the final container
FROM alpine:latest

# Set the working directory
WORKDIR /root/

RUN apk add curl
RUN curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s


# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Run the application
CMD ["./main"]
