
.PHONY: upgrade
upgrade:
	go get -u -v ./...

.PHONY: build
build:
	go build -v ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: test-cgo
test-cgo:
	go test -tags use_cgo,!cgo_ext -v ./...

.PHONY: test-cgo-ext
test-cgo-ext:
	go test -tags use_cgo,cgo_ext -v

.PHONY: test-nocgo
test-nocgo:
	go test -tags !use_cgo -v ./...