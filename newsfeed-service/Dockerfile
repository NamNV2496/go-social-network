# Build stage
FROM golang:1.22-alpine3.20 AS builder
WORKDIR /app/newsfeed-service

# Copy go.mod and go.sum to cache dependencies
COPY newsfeed-service/go.mod ./
COPY newsfeed-service/go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY newsfeed-service/ ./

# Build the binary
RUN go build -o /newsfeed cmd/main.go

# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /newsfeed .

EXPOSE 8082
CMD ["/app/newsfeed"]
