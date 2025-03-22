update:
	buf mod update ../api/

.PHONY: generate

generate:
	cp ./buf.gen.yaml ../api 
	buf generate ../api/post_core/v1/post.proto
	wire ./internal/wiring/ 
	rm ../api/buf.gen.yaml

run:
	go run cmd/main.go