BINARY_NAME=air
MAIN_PATH=./src/cmd/main.go
BUILD_DIR=./bin


.PHONY: all build run test clean lint

all: clean lint test build

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

run:
	go run $(MAIN_PATH)

test:
	go test -v -race -count=1 ./...

clean:
	rm -rf $(BUILD_DIR)

lint:
	golangci-lint run
