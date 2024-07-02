# Build stage
FROM golang:1.22-alpine3.20 AS builder
WORKDIR /app/user-service

# Copy go.mod and go.sum to cache dependencies
COPY user-service/go.mod ./
COPY user-service/go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY user-service/ ./

# Build the binary
RUN go build -o /user cmd/main.go

# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /user .

EXPOSE 8084
CMD ["/app/user"]
