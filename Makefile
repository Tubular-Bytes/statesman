VERSION := $(shell git describe --tags --abbrev=0)

test:
	go test ./... -cover

lint:
	golangci-lint run ./...

clean:
	rm -rf bin

build: amd64 arm64

amd64: ./bin/amd64/statesman-$(VERSION)
	GOARCH=amd64 GOOS=linux go build -o bin/amd64/statesman-$(VERSION) ./cmd/...

arm64: ./bin/arm64/statesman-$(VERSION)
	GOARCH=arm64 GOOS=linux go build -o bin/arm64/statesman-$(VERSION) ./cmd/...

build-images: image-amd64 image-arm64
	podman manifest create registry.0x42.in/terrence/statesman:$(VERSION)
	podman manifest add registry.0x42.in/terrence/statesman:$(VERSION) registry.0x42.in/terrence/statesman:$(VERSION)-amd64
	podman manifest add registry.0x42.in/terrence/statesman:$(VERSION) registry.0x42.in/terrence/statesman:$(VERSION)-arm64

push-images: build-images
	podman push registry.0x42.in/terrence/statesman:$(VERSION)-amd64
	podman push registry.0x42.in/terrence/statesman:$(VERSION)-arm64
	podman push registry.0x42.in/terrence/statesman:$(VERSION)

image-amd64: amd64
	cp bin/amd64/statesman-$(VERSION) statesman
	podman build --platform linux/amd64 -t registry.0x42.in/terrence/statesman:$(VERSION)-amd64 -f Dockerfile-release .
	rm statesman

image-arm64: arm64
	cp bin/arm64/statesman-$(VERSION) statesman
	podman build --platform linux/arm64 -t registry.0x42.in/terrence/statesman:$(VERSION)-arm64 -f Dockerfile-release .
	rm statesman