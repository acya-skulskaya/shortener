.PHONY: build-shortener
.PHONY: build-shortener-proto

ifdef VERSION
VERSION := $(VERSION)
else
VERSION := 0.0.1
endif

ifdef BUILD_DATE
BUILD_DATE := $(BUILD_DATE)
else
BUILD_DATE := $(shell date +'%Y-%m-%d_%H:%M:%S')
endif

ifdef COMMIT_HASH
COMMIT_HASH := $(COMMIT_HASH)
else
COMMIT_HASH := $(shell git rev-parse --short HEAD)
endif

build-shortener:
	go build -o ./cmd/shortener/shortener -ldflags "\
	-X main.buildVersion=$(VERSION) \
	-X main.buildDate=$(BUILD_DATE) \
	-X main.buildCommit=$(COMMIT_HASH)" \
	./cmd/shortener/

build-shortener-proto:
	protoc --go_out=.\
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
    --go_opt=default_api_level=API_OPAQUE \
    api/shortener/shortener.proto