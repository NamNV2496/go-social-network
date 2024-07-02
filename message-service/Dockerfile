# Build stage
FROM golang:1.22-alpine3.20 AS builder
WORKDIR /app/message-service

# Copy go.mod and go.sum to cache dependencies
COPY message-service/go.mod ./
COPY message-service/go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY message-service/ ./

# Build the binary
RUN go build -o /message cmd/main.go

# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /message .

EXPOSE 8081
CMD ["/app/message"]
