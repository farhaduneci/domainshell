APP_NAME=domainshell
BUILD_DIR=bin

.PHONY: all build build-all clean test install fmt vet lint

all: fmt vet build

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	# Requires golangci-lint to be installed
	golangci-lint run

build:
	go build -o $(APP_NAME) ./cmd/domainshell

build-all: clean
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 ./cmd/domainshell
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64 ./cmd/domainshell
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 ./cmd/domainshell
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 ./cmd/domainshell
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe ./cmd/domainshell
	@echo "Builds completed in $(BUILD_DIR)/"

test:
	go test ./...

clean:
	go clean
	rm -f $(APP_NAME)
	rm -rf $(BUILD_DIR)

install: build
	mv $(APP_NAME) $(GOPATH)/bin/$(APP_NAME)
