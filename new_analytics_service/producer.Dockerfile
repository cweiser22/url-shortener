# Use the official Golang image as a build stage
FROM golang:1.23-alpine

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
RUN go build -o main ./cmd/producer/main.go


RUN apk add curl
RUN curl -sSfL https://goblin.run/github.com/air-verse/air | sh


# Expose port 8002
EXPOSE 8002

# Run the application
CMD ["./main"]
