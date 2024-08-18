

.PHONY: upgrade
upgrade:
	go get -u -v ./...

.PHONY: build
build:
	go build -v ./...
