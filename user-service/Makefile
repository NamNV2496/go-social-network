update:
	buf mod update ../api/
generate:
	# buf generate ../api/
	buf generate ../api/user_core/v1/user.proto
	
wire:
	wire ./internal/wiring/
	
PORT ?= 5610

run:
	go run main.go --port $(PORT)
