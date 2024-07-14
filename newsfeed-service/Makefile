update:
	buf mod update ../api/

.PHONY: generate
generate:
	# buf generate ../api/user_core/v1/user.proto
	# buf generate ../api/newsfeed_core/v1/newsfeed.proto
	buf generate ../api/
	wire ./internal/wiring/

run:
	go run cmd/main.go