# Build stage
FROM golang:1.22-alpine3.20 AS builder
WORKDIR /app/http_gateway

# Copy go.mod and go.sum to cache dependencies
COPY http_gateway/go.mod ./
COPY http_gateway/go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY http_gateway/ ./

# Build the binary
RUN go build -o /http cmd/main.go

# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /http .

EXPOSE 8080
CMD ["/app/http"]
