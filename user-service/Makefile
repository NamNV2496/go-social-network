update:
	buf mod update ../api/
generate:
	# buf generate ../api/
	buf generate ../api/user_core/v1/user.proto
	
wire:
	wire ./internal/wiring/

run:
	go run cmd/main.go