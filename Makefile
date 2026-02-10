.PHONY: build test run clean

BINARY_NAME=tz-snipe
BUILD_DIR=build
CMD_PATH=cmd/tz-snipe/main.go

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)

test:
	go test ./...

# To pass arguments that start with -, use -- to separate make flags from app flags
# Example: make run -- --time 15:00
# Example: make run -- --tz +1000
run:
	@go run $(CMD_PATH) $(filter-out $@,$(MAKECMDGOALS))

clean:
	rm -rf $(BUILD_DIR)

# Catch-all rule to prevent make from complaining about unknown targets (the args)
%:
	@:
