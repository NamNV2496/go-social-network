# Build stage
FROM golang:1.22-alpine3.20 AS builder
WORKDIR /app/post-service

# Copy go.mod and go.sum to cache dependencies
COPY post-service/go.mod ./
COPY post-service/go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY post-service/ ./

# Build the binary
RUN go build -o /post cmd/main.go

# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /post .

EXPOSE 8083
CMD ["/app/post"]
