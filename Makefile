export node_bins=$(PWD)/node_modules/.bin
export uglifyjs=$(node_bins)/uglifyjs
export gulp=$(node_bins)/gulp
export webpack=$(node_bins)/webpack
export GO111MODULE=on

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	ROCKSDB_CFLAGS := $(shell pkg-config --cflags-only-I liblz4 libzstd) -I/opt/homebrew/opt/snappy/include
	ROCKSDB_LDFLAGS := $(shell pkg-config --libs-only-L liblz4 libzstd) -L/opt/homebrew/opt/snappy/lib
	WEBP_CFLAGS = $(shell pkg-config --cflags libwebp)
    WEBP_LDFLAGS = $(shell pkg-config --libs-only-L libwebp)
    GO_BUILD_TAGS = -tags "libsqlite3"
endif

ifeq ($(UNAME_S),Linux)
	ROCKSDB_LDFLAGS := -lbz2
    GO_BUILD_TAGS = -tags "libsqlite3 linux"
endif

.PHONY: client server imager test

all: client server

client: client_deps proto_client
	node esbuild.config.js --css --js

client_deps:
	npm install --include=dev --progress false --depth 0

css:
	node esbuild.config.js --css

js:
	node esbuild.config.js --js

proto: proto_client proto_server

proto_client:
	npx protoc --ts_out client/typings --proto_path pb --experimental_allow_proto3_optional pb/nekotv.proto pb/posts.proto

proto_server:
	protoc --go_out=pb --proto_path=pb --go_opt=paths=source_relative --experimental_allow_proto3_optional pb/*.proto

generate:
	go generate .


server: proto_server
	go generate
	CGO_CFLAGS="$(ROCKSDB_CFLAGS) $(WEBP_CFLAGS)" CGO_LDFLAGS="$(ROCKSDB_LDFLAGS) $(WEBP_LDFLAGS)" go build -v $(GO_BUILD_TAGS)

client_clean:
	rm -rf www/js www/css/*.css www/css/maps node_modules manifest.json

clean: client_clean
	rm -rf .build .ffmpeg .package target meguca-*.zip meguca-*.tar.xz meguca meguca.exe server/pkg

test:
	go test --race ./...

test_no_race:
	go test ./...

test_docker:
	docker-compose build
	docker-compose run --rm -e CI=true meguca make test
