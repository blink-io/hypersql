
.PHONY: upgrade
upgrade:
	go get -u -v ./...

.PHONY: build
build:
	go build -v ./...

.PHONY: build-with-sqlite
build-with-sqlite:
	go build -tags sqlite -v ./...

.PHONY: build-with-sqlite-cgo
build-with-sqlite-cgo:
	go build -tags sqlite,sqlite_cgo -v ./...

.PHONY: build-all
build-all: build build-with-sqlite build-with-sqlite-cgo

.PHONY: test
test:
	go test -v ./...

.PHONY: test-sqlite
test-sqlite:
	go test -tags sqlite -v ./...

.PHONY: test-sqlite-cgo
test-sqlite-cgo:
	go test -tags sqlite,sqlite_cgo -v ./...