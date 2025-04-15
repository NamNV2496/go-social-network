
generate:
	buf dep update ../api/
	buf generate ../api/
	wire ./internal/wiring/

run:
	go run cmd/main.go