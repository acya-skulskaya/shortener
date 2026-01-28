.PHONY: build-shortener

VERSION := 0.0.1
BUILD_DATE := $(shell date +'%Y-%m-%d_%H:%M:%S')
COMMIT_HASH := $(shell git rev-parse --short HEAD)

build-shortener:
	go build -o ./cmd/shortener/shortener -ldflags "\
	-X main.buildVersion=$(VERSION) \
	-X main.buildDate=$(BUILD_DATE) \
	-X main.buildCommit=$(COMMIT_HASH)" \
	./cmd/shortener/