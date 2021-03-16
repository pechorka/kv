test:
	go test github.com/pechorka/kv/internal/store
	go test github.com/pechorka/kv/app/kv-api/handlers

build:
	docker build \
		-f docker/dockerfile.kv-api \
		-t kv-api-amd64:1.0 \
		.


run:
	docker run -d -p 8080:8080 kv-api-amd64:1.0