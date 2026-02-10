.PHONY: build test run clean

BINARY_NAME=tz-snipe
BUILD_DIR=build
CMD_PATH=cmd/tz-snipe/main.go

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)

test:
	go test ./...

run:
	go run $(CMD_PATH) $(args)

clean:
	rm -rf $(BUILD_DIR)
