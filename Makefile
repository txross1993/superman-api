.PHONY: all
all: build-api

build-api-mac:
	CGO_ENABLED=1 GOOS=darwin CGO_LDFLAGS="-g -O2 -L/usr/local/opt/openssl/lib" go build -a -installsuffix cgo -ldflags "-linkmode external -extldflags -static" -o app main.go
	chmod +x app

build-api:
	CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags "-linkmode external -extldflags -static" -o app main.go
	chmod +x app

docker-run:
	docker-compose -f build/package/docker-compose.yaml up --build

docker-build:
	docker-compose -f build/package/docker-compose.yaml build --force-rm

docker-test:
	docker-compose -f build/test/docker-compose.yaml up --exit-code-from superman-api-test

init-db:
	mkdir -p /local-db
	touch /local-db/local.db
	chmod -R 644 /local-db/local.db


lint:
	golint -set_exit_status $(shell go list ./... | grep -v /vendor/)

test:
	go test $(shell go list ./... | grep -v /vendor/)

clean:
	go clean -cache
	go clean -testcache

get-dependencies:
	go mod download
