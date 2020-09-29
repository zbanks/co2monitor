GO=/usr/lib/go-1.10/bin/go

install:
	$(GO) install -v

build:
	$(GO) build -v

build_arm:
	GOOS=linux GOARCH=arm GOARM=6 $(GO) build -o co2monitor.armv6
	GOOS=linux GOARCH=arm GOARM=7 $(GO) build -o co2monitor.armv7

test:
	$(GO) test --race -v ./...

.PHONY: install build build_arm test
