update:
	buf mod update ../api/

.PHONY: generate
generate:
	buf generate ../api/post_core/v1/post.proto
	wire ./internal/wiring/

run:
	go run cmd/main.go