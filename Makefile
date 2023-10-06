BINARY_NAME=bth-speaker-on
UP_INTERVAL=60

build: export GOOS=darwin
build: export GOARCH=arm64
build:
	@go build -o bin/$(BINARY_NAME) -v

run: build
	@./bin/$(BINARY_NAME) start up-interval=$(UP_INTERVAL)