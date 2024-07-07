# Makefile for Golang project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOFMT=$(GOCMD) fmt
BINARY_NAME=tcp-hb
BINARY_UNIX=$(BINARY_NAME)_linux
BINARY_MAC=$(BINARY_NAME)_mac

# Directories
SRC_DIR=.
INTERNAL_DIR=./internal

# Default target executed when no arguments are given to make.
default: build

# Build the project
build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(SRC_DIR)/main.go

# Run the project
run:
	$(GOBUILD) -o $(BINARY_NAME) -v $(SRC_DIR)/main.go
	./$(BINARY_NAME)

# Test the project
test:
	$(GOTEST) -v ./...

# Clean the project
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_MAC)

# Format the project
fmt:
	$(GOFMT) ./...

# Cross compile for Linux
build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v $(SRC_DIR)/main.go

# Cross compile for MacOS
build-mac:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_MAC) -v $(SRC_DIR)/main.go

.PHONY: build run test clean fmt build-linux build-mac
