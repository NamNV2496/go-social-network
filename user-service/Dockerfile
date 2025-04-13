# Build stage
FROM golang:1.22-alpine3.20 AS builder
WORKDIR /app

# Copy go.mod and go.sum to cache dependencies
COPY user-service/go.mod user-service/go.sum ./user-service/
WORKDIR /app/user-service
RUN go mod download

# Copy the rest of the source code
COPY user-service/ ./

# Build the binary
RUN go build -o /user-service main.go

# Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /user-service .

# Copy the templates directory into the Docker image
COPY user-service/templates /app/templates

EXPOSE 8084
CMD ["/app/user-service"]

